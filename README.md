# Simplebackup EC2

## Description

To simplify the Amazon EBS snapshot creation and deletion management.

## Usage

- Create a snapshot of the ebs with instances.

```
c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
s, _ := simplebackupec2.NewService(c)
err := s.CreateSnapshots("i-xxxxxxxx")
```

- Deleted until a specified number of a snapshot of the ebs with instances.

```
c := simplebackupec2.NewConfig().WithRegion("ap-northeast-1").WithCredentials(creds)
s, _ := simplebackupec2.NewService(c)
err := s.RotateSnapshots("i-xxxxxxxx", 5)
```

- Required IAM Policy.

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:CreateSnapshot",
                "ec2:CreateTags",
                "ec2:DeleteSnapshot",
                "ec2:DescribeSnapshots",
                "ec2:DescribeInstances",
                "ec2:CreateImage",
                "ec2:DeregisterImage",
                "ec2:DescribeImages"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```


## Install

To install, use `go get`:

```bash
$ go get -d github.com/simplebackup/ec2
$ cd $GOPATH/src/github.com/simplebackup/ec2
```

## Warning!

Test does not exist yet.

## Contribution

1. Fork ([https://github.com/simplebackup/ec2/fork](https://github.com/simplebackup/ec2/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[youyo](https://github.com/youyo)
