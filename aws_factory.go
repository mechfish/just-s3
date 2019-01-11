package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type AwsFactory interface {
	NewS3(...*aws.Config) s3iface.S3API
}

type DefaultAwsFactory struct {
	Session *session.Session
}

func NewAwsFactory(cfgs ...*aws.Config) AwsFactory {
	return &DefaultAwsFactory{
		Session: session.Must(session.NewSession(cfgs...)),
	}
}

func (f *DefaultAwsFactory) NewS3(cfgs ...*aws.Config) s3iface.S3API {
	return s3.New(f.Session, cfgs...)
}
