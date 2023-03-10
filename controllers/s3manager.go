/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"crypto/x509"
	"net/http"

	s3v1 "github.com/giubacc/s3gw-operator/api/v1"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

type Manager struct {
	s3Client          *minio.Client
	connectionDetails ConnectionDetails
}

type ConnectionDetails struct {
	Endpoint        string
	UseSSL          bool
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	CA              []byte
}

// Validate makes sure the provided S3 settings are valid
func (details *ConnectionDetails) ValidateS3ConnectionDetails() error {
	if details.AccessKeyID == "" ||
		details.SecretAccessKey == "" ||
		details.Endpoint == "" {
		return errors.New("invalid S3ConnectionDetails")
	}
	return nil
}

// New returns an instance of an s3 manager
func NewS3Manager(connectionDetails ConnectionDetails) (*Manager, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	opts := &minio.Options{
		Creds:     credentials.NewStaticV4(connectionDetails.AccessKeyID, connectionDetails.SecretAccessKey, ""),
		Secure:    connectionDetails.UseSSL,
		Transport: transport,
		Region:    connectionDetails.Region,
	}

	if len(connectionDetails.CA) > 0 {
		rootCAs := x509.NewCertPool()
		if ok := rootCAs.AppendCertsFromPEM(connectionDetails.CA); !ok {
			return nil, errors.New("error with connectionDetails.CA")
		}

		tlsConfig := transport.TLSClientConfig.Clone()
		tlsConfig.RootCAs = rootCAs

		opts.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	}

	minioClient, err := minio.New(
		connectionDetails.Endpoint,
		opts,
	)
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		connectionDetails: connectionDetails,
		s3Client:          minioClient,
	}

	return manager, nil
}

// GetConnectionDetails retrieves s3 details
func GetS3ConnectionDetails(AccessKeyID string, SecretAccessKey string, Region string, Endpoint string, UseSSL bool) (ConnectionDetails, error) {
	details := ConnectionDetails{}

	details.AccessKeyID = AccessKeyID
	details.SecretAccessKey = SecretAccessKey
	details.Region = Region
	details.Endpoint = Endpoint
	details.UseSSL = UseSSL

	return details, nil
}

// EnsureBucketCreated creates a bucket
func (m *Manager) EnsureBucketCreated(ctx context.Context, bucket *s3v1.Bucket) error {
	exists, err := m.s3Client.BucketExists(ctx, bucket.Spec.Name)
	if err != nil {
		DebugLogger.Errorf("BucketExists:%s", err.Error())
		return err
	}

	if !exists {
		if err := m.s3Client.MakeBucket(ctx, bucket.Spec.Name,
			minio.MakeBucketOptions{Region: m.connectionDetails.Region, ObjectLocking: bucket.Spec.ObjectLocking}); err != nil {
			DebugLogger.Errorf("MakeBucket:%s", err.Error())
			return err
		}
	}

	if bucket.Spec.Versioning {
		bconf, err := m.s3Client.GetBucketVersioning(ctx, bucket.Spec.Name)
		if err != nil {
			DebugLogger.Errorf("GetBucketVersioning:%s", err.Error())
			return err
		}
		if !bconf.Enabled() {
			if err := m.s3Client.EnableVersioning(ctx, bucket.Spec.Name); err != nil {
				DebugLogger.Errorf("EnableVersioning:%s", err.Error())
				return err
			}
		}
	}

	return nil
}

// EnsureBucketDeleted delete a bucket
func (m *Manager) EnsureBucketDeleted(ctx context.Context, bucket *s3v1.Bucket) error {
	var retError error = nil

	exists, err := m.s3Client.BucketExists(ctx, bucket.Spec.Name)
	if err != nil {
		DebugLogger.Errorf("BucketExists:%s", err.Error())
		return err
	}

	if exists {
		objectsCh := make(chan minio.ObjectInfo)

		// Send object names that are needed to be removed to objectsCh
		go func() {
			defer close(objectsCh)
			// List all objects from a bucket-name with a matching prefix.
			opts := minio.ListObjectsOptions{Prefix: "", Recursive: true}
			for object := range m.s3Client.ListObjects(context.Background(), bucket.Spec.Name, opts) {
				if object.Err != nil {
					DebugLogger.Errorf("Failed to list objects, error: %s", object.Err.Error())
				}
				objectsCh <- object
			}
		}()

		// Call RemoveObjects API
		errorCh := m.s3Client.RemoveObjects(context.Background(), bucket.Spec.Name, objectsCh, minio.RemoveObjectsOptions{})

		// Print errors received from RemoveObjects API
		for e := range errorCh {
			DebugLogger.Errorf("Failed to remove: %s, error: %s", e.ObjectName, e.Err.Error())
			//save first error received
			if retError != nil {
				retError = e.Err
			}
		}

		if retError == nil {
			return m.s3Client.RemoveBucket(ctx, bucket.Spec.Name)
		}
	}

	return retError
}
