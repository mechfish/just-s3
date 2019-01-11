package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// parseFullS3Url tries to parse the given string as an `s3://`
// url. If it succeeds, it returns the S3 bucket and object names
// along with `true`. Otherwise it returns empty strings and `false`.
func parseS3Url(url string) (string, string, bool) {
	re := regexp.MustCompile(`\As3://([^/]+)/([^/].*)`)
	m := re.FindStringSubmatch(url)
	if m == nil {
		return "", "", false
	}
	return m[1], m[2], true
}

func urlReader(url string, awsFactory AwsFactory) (io.ReadCloser, error) {
	b, o, isS3 := parseS3Url(url)
	if isS3 {
		result, err := awsFactory.NewS3().GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(b),
				Key:    aws.String(o),
			})
		if err != nil {
			return nil, err
		}
		return result.Body, nil
	}
	return os.Open(url)
}

func writeUrl(url string, data io.Reader, awsFactory AwsFactory) error {
	var targetFile *os.File
	var err error
	b, o, isS3 := parseS3Url(url)
	if isS3 {
		if targetFile, err = ioutil.TempFile(``, `just-s3-cp`); err != nil {
			return err
		}
		defer os.Remove(targetFile.Name())
	} else {
		if targetFile, err = os.Create(url); err != nil {
			return err
		}
	}
	if _, err = io.Copy(targetFile, data); err != nil {
		return err
	}
	if !isS3 {
		return nil
	}
	if _, err = targetFile.Seek(0, 0); err != nil {
		return err
	}
	_, err = awsFactory.NewS3().PutObject(&s3.PutObjectInput{
		Bucket: aws.String(b),
		Key:    aws.String(o),
		Body:   targetFile,
	})
	return err
}

func copyObject(srcPath string, destPath string, awsFactory AwsFactory) error {
	fmt.Printf("Copying %s to %s\n", srcPath, destPath)
	reader, err := urlReader(srcPath, awsFactory)
	if err != nil {
		return err
	}
	return writeUrl(destPath, reader, awsFactory)
}
