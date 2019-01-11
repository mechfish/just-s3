package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// An AwsFactory wraps connections to the AWS API.
//
// Unit tests may use a fake factory that returns fake connections.
type AwsFactory interface {
	NewS3(...*aws.Config) s3iface.S3API
}

// DefaultAwsFactory returns real connections to the AWS API.
type DefaultAwsFactory struct {
	Session *session.Session
}

// NewAwsFactory accepts the same arguments as AWS's
// `session.NewSession` and returns an AwsFactory that makes real
// connections to AWS.
func NewAwsFactory(cfgs ...*aws.Config) AwsFactory {
	return &DefaultAwsFactory{
		Session: session.Must(session.NewSession(cfgs...)),
	}
}

// NewS3 return an AWS S3 client.
func (f *DefaultAwsFactory) NewS3(cfgs ...*aws.Config) s3iface.S3API {
	return s3.New(f.Session, cfgs...)
}
