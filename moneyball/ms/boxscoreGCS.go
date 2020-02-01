package ms

/**
Copyright (c) 2020 DXC Technology - Dan Hushon. All rights reserved

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc., DXC Technology nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
)

//sorting storage metadata attributes
/// 3 classes of attributes we need to manage
// 1 bucket
// 2 default objects in bucket
// 3 object

func getCreds() option.ClientOption {
	creds, exists := os.LookupEnv("BOXSCORE_CREDS")
	if !(exists) {
		return nil
	}
	return option.WithCredentialsFile(creds)
}

//Write write the byte.Buffer to the named object->bucket, inclusive of a set of attributes
func Write(b *bytes.Buffer, projectID *string, bucketName string, objectName string, attr map[string]interface{}) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, getCreds())
	if e, ok := err.(*googleapi.Error); ok {
		if e.Code == 409 {
			log.Printf("error %v\n", e.Error())
		}
	}
	// create bucket handle
	bkt := client.Bucket(bucketName)

	//test if bkthandle exists?
	if err := existsBucket(ctx, bkt); err != nil {
		bkt, err = createBucket(ctx, projectID, bucketName)
		if err != nil {
			return err
		}
	}

	// TODO: sort attributes on the bucket
	ba, err := bkt.Attrs(ctx)
	boa, err := bkt.DefaultObjectACL().List(ctx)
	if err != nil {
		return err
	}
	log.Printf("Attributes %v\n", ba)
	log.Printf("Default Object ACL %v\n", boa)
	// if not sync'd need to error so as not to create insecure stuff

	wc := bkt.Object(objectName).NewWriter(ctx)
	wc.ContentType = "text/json"

	if _, err = io.Copy(wc, b); err != nil {
		log.Printf("Error: %v\n", err)
	}

	//push new attr triples at some point in future
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	fmt.Println("updated object:", wc.Attrs())

	if err := wc.Close(); err != nil {
		log.Printf("Error: %v\n", err)
	}
	return nil

}

func existsBucket(ctx context.Context, bktHandle *storage.BucketHandle) error {
	//test if bkthandle exists?
	_, err := bktHandle.Attrs(ctx)
	return err
}

func createBucket(ctx context.Context, projectID *string, bucketName string) (*storage.BucketHandle, error) {
	// Creates a client.
	client, err := storage.NewClient(ctx, getCreds())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	// Creates the new bucket, throws error if create fails
	if err := bucket.Create(ctx, *projectID, nil); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == 409 {
				log.Printf("error - bucket name must be globally unique pls try again%v\n", e.Error())
			} else {
				log.Fatalf("Failed to create bucket: %v", e)
			}
			return nil, e
		}
		return nil, err
	}
	log.Printf("Bucket %v created.\n", bucketName)
	return bucket, nil
}

func deleteBucketByName(ctx context.Context, projectID *string, bucketName string) error {
	// Creates a client.
	client, err := storage.NewClient(ctx, getCreds())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return err
	}
	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	return deleteBucket(ctx, bucket)
}

func deleteBucket(ctx context.Context, bkthandle *storage.BucketHandle) error {
	//teardown - bucket deleted
	err := bkthandle.Delete(ctx)
	return err
}
