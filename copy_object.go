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
// URL. If it succeeds, it returns the S3 bucket and object names
// along with `true`. Otherwise it returns empty strings and `false`.
func parseS3Url(url string) (string, string, bool) {
	re := regexp.MustCompile(`\As3://([^/]+)/([^/].*)`)
	m := re.FindStringSubmatch(url)
	if m == nil {
		return "", "", false
	}
	return m[1], m[2], true
}

// urlReader returns a ReadCloser that will read from the given URL,
// whether that URL points to s3 or is a path to a file on the local
// disk.
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

// writeUrl reads data from the given reader and writes it to the
// given URL, whether that URL points to s3 or to a file on local
// disk.
func writeUrl(url string, reader io.Reader, awsFactory AwsFactory) error {
	var targetFile *os.File
	var err error
	b, o, isS3 := parseS3Url(url)
	if isS3 {
		// To write to S3 we need an io.ReadSeeker, not a plain
		// io.Reader (which does not support Seek()); this is probably
		// because the S3 Put transaction needs to be signed, so it
		// needs to scan the content twice: once to generate the
		// signature, and once to send it over the network.
		//
		// So we copy the contents of the reader to a temporary file
		// on disk, which will support Seek(), and then upload the
		// temporary file to S3.
		if targetFile, err = ioutil.TempFile(``, `just-s3-cp`); err != nil {
			return err
		}
		defer os.Remove(targetFile.Name())
	} else {
		if targetFile, err = os.Create(url); err != nil {
			return err
		}
	}
	// Copy from the reader to the file on disk.
	if _, err = io.Copy(targetFile, reader); err != nil {
		return err
	}
	if !isS3 {
		// The file on disk is the final destination; we're done
		return nil
	}
	// Copy the temporary file to S3
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

// copyObject copies the contents of the srcPath to the
// destPath. These paths may be S3 URLs, or paths to files on the
// local disk.
func copyObject(srcPath string, destPath string, awsFactory AwsFactory) error {
	fmt.Printf("Copying %s to %s\n", srcPath, destPath)
	reader, err := urlReader(srcPath, awsFactory)
	if err != nil {
		return err
	}
	defer reader.Close()
	return writeUrl(destPath, reader, awsFactory)
}
