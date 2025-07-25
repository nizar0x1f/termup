package s3storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cheggaaa/pb/v3"
	"github.com/nizar0x1f/termup/pkg/config"
)

type UploadOptions struct {
	InsecureTLS bool
}

type ProgressCallback func(uploaded int64)

type progressReader struct {
	reader   io.ReadSeeker
	total    int64
	read     int64
	callback ProgressCallback
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.read += int64(n)
	if pr.callback != nil {
		pr.callback(pr.read)
	}
	return n, err
}

func (pr *progressReader) Seek(offset int64, whence int) (int64, error) {
	pos, err := pr.reader.Seek(offset, whence)
	if err == nil {
		pr.read = pos
		if pr.callback != nil {
			pr.callback(pr.read)
		}
	}
	return pos, err
}

func Upload(cfg *config.Config, filePath string) (string, error) {
	return UploadWithOptions(cfg, filePath, nil)
}

func UploadWithProgress(cfg *config.Config, filePath string, progressCallback ProgressCallback) (string, error) {
	return UploadWithOptionsAndProgress(cfg, filePath, nil, progressCallback)
}

func UploadWithOptions(cfg *config.Config, filePath string, opts *UploadOptions) (string, error) {
	return UploadWithOptionsAndProgress(cfg, filePath, opts, nil)
}

func UploadWithOptionsAndProgress(cfg *config.Config, filePath string, opts *UploadOptions, progressCallback ProgressCallback) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	var httpClient *http.Client
	if opts != nil && opts.InsecureTLS {
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	configOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
		awsconfig.WithRegion("auto"),
	}

	if httpClient != nil {
		configOptions = append(configOptions, awsconfig.WithHTTPClient(httpClient))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return "", fmt.Errorf("failed to load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true

		o.BaseEndpoint = aws.String(cfg.Endpoint)
	})

	fileName := filepath.Base(filePath)

	var body io.ReadSeeker = file

	if progressCallback != nil {
		body = &progressReader{
			reader:   file,
			total:    fileInfo.Size(),
			callback: progressCallback,
		}
	} else {

		bar := pb.Full.Start64(fileInfo.Size())
		defer bar.Finish()
		defer bar.SetCurrent(fileInfo.Size())
	}

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(fileName),
		Body:   body,
	})

	if err != nil {

		return "", fmt.Errorf("failed to upload file '%s' to bucket '%s': %w", fileName, cfg.Bucket, err)
	}

	return fmt.Sprintf("%s/%s", strings.TrimSuffix(cfg.PublicUrl, "/"), fileName), nil
}
