package simplebackupec2

import (
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

const description = "Created by simplebackup/ec2 from"

// Service is wrapping configuration *ec2.EC2
type Service struct {
	*ec2.EC2
}

type snapshot struct {
	ID        string
	StartTime int64
}

type snapshots []snapshot

// NewConfig returns a new Config pointer that can be chained with builder methods to set multiple configuration values inline without using pointers.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
func NewConfig() *aws.Config {
	return &aws.Config{}
}

// NewService creates a new instance of the EC2 client with a session.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, err := simplebackupec2.NewService(c)
func NewService(c *aws.Config) (*Service, error) {
	sess, err := session.NewSession(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new session")
	}
	return &Service{ec2.New(sess)}, nil
}

// CreateSnapshots create EBS snapshot.
// To create a snapshot of all the volumes that have instances.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  err := s.CreateSnapshots("i-xxxxxxxx")
func (s *Service) CreateSnapshots(instanceID string) error {
	tag, err := s.readNameTag(instanceID)
	if err != nil {
		return errors.Wrap(err, "failed to read name tag")
	}
	volumeIDs, err := s.describeAllVolumeIDs(instanceID)
	if err != nil {
		return errors.Wrap(err, "failed to describe all volumes")
	}
	for _, volumeID := range volumeIDs {
		snapshotID, err := s.createSnapshot(volumeID)
		if err != nil {
			return errors.Wrap(err, "failed to create snapshots")
		}
		if err := s.setNameTagToSnapshot(snapshotID, tag); err != nil {
			return errors.Wrap(err, "failed to set name tag to snapshot")
		}
	}
	return nil
}

// RotateSnapshot manages the number of snapshot of a specific volume.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  err := s.RotateSnapshot("v-xxxxxxxx", 5)
func (s *Service) RotateSnapshot(volumeID string, i int) error {
	snapshots, err := s.describeSnapshots(volumeID)
	if err != nil {
		return errors.Wrap(err, "failed to describe snapshot.")
	}
	snapshotIDs := sortSnapshots(snapshots)
	for len(snapshotIDs) > i {
		params := &ec2.DeleteSnapshotInput{
			SnapshotId: aws.String(snapshotIDs[0].ID),
		}
		_, err := s.DeleteSnapshot(params)
		if err != nil {
			return errors.Wrap(err, "failed to delete snapshot.")
		}
		snapshotIDs = append(snapshotIDs[1:])
	}
	return nil
}

// RotateSnapshots manages the number of snapshot of a specific instance's volumes.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  err := s.RotateSnapshots("i-xxxxxxxx", 5)
func (s *Service) RotateSnapshots(instanceID string, i int) error {
	volumeIDs, err := s.describeAllVolumeIDs(instanceID)
	if err != nil {
		return errors.Wrap(err, "failed to describe all volumes")
	}
	for _, volumeID := range volumeIDs {
		if err = s.RotateSnapshot(volumeID, i); err != nil {
			return errors.Wrap(err, "failed to delete snapshot: "+volumeID)
		}
	}
	return nil
}

// RegisterAMI create New AMI.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  imageID, err := s.RegisterAMI("i-xxxxxxxx", true)
func (s *Service) RegisterAMI(instanceID string, noReboot bool) (string, error) {
	tag, err := s.readNameTag(instanceID)
	if err != nil {
		return "", errors.Wrap(err, "failed to read name tag")
	}
	d := description + " " + instanceID +
		" at " + strconv.FormatInt(time.Now().Unix(), 10)
	params := &ec2.CreateImageInput{
		InstanceId:  aws.String(instanceID),
		Name:        aws.String(d),
		Description: aws.String(d),
		NoReboot:    aws.Bool(noReboot),
	}
	resp, err := s.CreateImage(params)
	if err != nil {
		return "", errors.Wrap(err, "failed to create AMI")
	}
	if err := s.setNameTagToAMI(*resp.ImageId, tag); err != nil {
		return *resp.ImageId, errors.Wrap(err, "create AMI is successful, but failed to set name tag")
	}
	return *resp.ImageId, nil
}

// DeregisterAMI deregister AMI.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  err := s.DeregisterAMI("ami-xxxxxxxx")
func (s *Service) DeregisterAMI(imageID string) error {
	params := &ec2.DeregisterImageInput{
		ImageId: aws.String(imageID),
	}
	_, err := s.DeregisterImage(params)
	if err != nil {
		return errors.Wrap(err, "faild to deregister AMI")
	}
	return nil
}

func (s *Service) describeInstances(instanceID string) (*ec2.DescribeInstancesOutput, error) {
	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}
	resp, err := s.DescribeInstances(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to describeInstances")
	}
	return resp, nil
}

func (s *Service) readNameTag(instanceID string) (string, error) {
	resp, err := s.describeInstances(instanceID)
	if err != nil {
		return "", err
	}
	tag, err := func(resp *ec2.DescribeInstancesOutput) (string, error) {
		for _, res := range resp.Reservations {
			for _, res := range res.Instances {
				for _, res := range res.Tags {
					if *res.Key == "Name" {
						return *res.Value, nil
					}
				}
			}
		}
		return "", nil
	}(resp)
	if err != nil {
		return "", errors.Wrap(err, "failed to read name tag")
	}
	return tag, nil
}

func (s *Service) describeAllVolumeIDs(instanceID string) ([]string, error) {
	resp, err := s.describeInstances(instanceID)
	if err != nil {
		return nil, err
	}
	var v []string
	for _, resp := range resp.Reservations {
		for _, resp := range resp.Instances {
			for _, resp := range resp.BlockDeviceMappings {
				v = append(v, *resp.Ebs.VolumeId)
			}
		}
	}
	return v, nil
}

func (s *Service) setNameTagToSnapshot(snapshotID, tag string) error {
	params := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(snapshotID),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(tag),
			},
		},
	}
	_, err := s.CreateTags(params)
	return err
}

func (s *Service) setNameTagToAMI(imageID, tag string) error {
	params := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(imageID),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(tag),
			},
		},
	}
	_, err := s.CreateTags(params)
	return err
}

func (s *Service) createSnapshot(volumeID string) (string, error) {
	d := description + " " + volumeID
	params := &ec2.CreateSnapshotInput{
		VolumeId:    aws.String(volumeID),
		Description: aws.String(d),
	}
	resp, err := s.CreateSnapshot(params)
	if err != nil {
		return "", errors.Wrap(err, "failed to create snapshot")
	}
	return *resp.SnapshotId, nil

}

func (s *Service) describeSnapshots(volumeID string) (*ec2.DescribeSnapshotsOutput, error) {
	params := &ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					aws.String(volumeID),
				},
			},
		},
	}
	return s.DescribeSnapshots(params)
}

func isOwn(d string) bool {
	if m, _ := regexp.MatchString(description+".*", d); !m {
		return false
	}
	return true
}

func (p snapshots) Len() int {
	return len(p)
}

func (p snapshots) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p snapshots) Less(i, j int) bool {
	return p[i].StartTime < p[j].StartTime
}

func sortSnapshots(s *ec2.DescribeSnapshotsOutput) snapshots {
	var snapshotIDs snapshots = make([]snapshot, 1)
	for _, res := range s.Snapshots {
		if snapshotIDs[0].StartTime == 0 {
			if isOwn(*res.Description) {
				snapshotIDs[0].ID = *res.SnapshotId
				snapshotIDs[0].StartTime = res.StartTime.Unix()
			}
		} else {
			if isOwn(*res.Description) {
				snapshotIDs = append(snapshotIDs,
					snapshot{*res.SnapshotId, res.StartTime.Unix()},
				)
			}
		}
	}
	sort.Sort(snapshotIDs)
	return snapshotIDs
}
