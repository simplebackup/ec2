package simplebackupec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

// Service is wrapping configuration *ec2.EC2
type Service struct {
	*ec2.EC2
}

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
func (s *Service) CreateSnapshots(instanceID string) (string, error) {
	pp.Print(s)
	pp.Print(instanceID)
	return instanceID, nil
}

// RotateSnapshot manages the number of snapshot of a specific volume.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  err := s.RotateSnapshot("v-xxxxxxxx", 5)
func (s *Service) RotateSnapshot(volumeID string, i int) error {
	pp.Print(s)
	pp.Print(volumeID)
	pp.Print(i)
	return nil
}

// RotateSnapshots manages the number of snapshot of a specific instance's volumes.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  err := s.RotateSnapshots("i-xxxxxxxx", 5)
func (s *Service) RotateSnapshots(instanceID string, i int) error {
	pp.Print(s)
	pp.Print(instanceID)
	pp.Print(i)
	return nil
}

// RegisterAMI create New AMI.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  imageID, err := s.RegisterAMI("i-xxxxxxxx")
func (s *Service) RegisterAMI(instanceID string) (string, error) {
	pp.Print(s)
	pp.Print(instanceID)
	return instanceID, nil
}

// DeregisterAMI deregister AMI.
//  c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
//  s, _ := simplebackupec2.NewService(c)
//  err := s.DeregisterAMI("ami-xxxxxxxx")
func (s *Service) DeregisterAMI(imageID string) error {
	pp.Print(s)
	pp.Print(imageID)
	return nil
}

func main() {
	c := NewConfig().WithRegion("us-west-2").WithMaxRetries(10)
	s, err := NewService(c)
	pp.Print(s)
	pp.Print(err)
}
