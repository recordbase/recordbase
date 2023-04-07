/*
 * Copyright (c) 2022-2023 Zander Schwid & Co. LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package recordbase

import (
	"context"
	"github.com/codeallergy/glue"
	"github.com/recordbase/recordpb"
	"google.golang.org/grpc"
	"reflect"
)

const (
	ServiceName = "recordbase.RecordService"
)

var ClientClass = reflect.TypeOf((*Client)(nil)).Elem()

type RecordEntryEvent struct {
	Entry      *recordpb.RecordEntry
	Err        error
}

type MapEntryEvent struct {
	Entry      *recordpb.MapEntry
	Err        error
}

type FileContentEvent struct {
	Content    *recordpb.FileContent
	Err        error
}

type Client interface {
	glue.DisposableBean

	//
	// Gets metadata about using attributes
	//
	GetCounts(ctx context.Context, in *recordpb.TenantRequest) (*recordpb.Counts, error)

	//
	// Quick user lookup request
	//
	Lookup(ctx context.Context, in *recordpb.LookupRequest) (*recordpb.RecordEntry, error)

	//
	// Search users by indexed non-unique attributes
	//
	Search(ctx context.Context, in *recordpb.SearchRequest) (entries <- chan RecordEntryEvent, cancel func(), err error)

	//
	// Get user with all attributes
	//
	Get(ctx context.Context, in *recordpb.GetRequest) (*recordpb.RecordEntry, error)

	//
	// Create user, returns new user_id
	//
	Create(ctx context.Context, in *recordpb.CreateRequest) (*recordpb.CreateResponse, error)

	//
	// Delete user request (sets TTL to all PII data for particular user)
	//
	Delete(ctx context.Context, in *recordpb.DeleteRequest) error

	//
	// Update user attributes
	//
	Update(ctx context.Context, in *recordpb.UpdateRequest) error

	//
	// Upload File
	//
	UploadFile(c context.Context) (sink chan <- *recordpb.UploadFileContent, _ <- chan error)

	//
	// Download File
	//
	DownloadFile(ctx context.Context, in *recordpb.DownloadFileRequest, opts ...grpc.CallOption)  (entries <- chan FileContentEvent, cancel func(), err error)

	//
	// Scan users
	//
	Scan(ctx context.Context, in *recordpb.ScanRequest) (entries <- chan RecordEntryEvent, cancel func(), err error)

	//
	// Allocate user id range
	//
	AddKeyRange(ctx context.Context, in *recordpb.KeyRange) error

	//
	// Gets user id ranges and etc
	//
	GetKeyCapacity(ctx context.Context, in *recordpb.TenantRequest) (*recordpb.KeyCapacity, error)

	//
	// Get map value associated with the record
	//
	MapGet(ctx context.Context, in *recordpb.MapGetRequest) (*recordpb.MapEntry, error)

	//
	// Put map value associated with the record. Returns old value.
	//
	MapPut(ctx context.Context, in *recordpb.MapPutRequest) (*recordpb.MapValue, error)

	//
	// Remove map value associated with the record. Returns old value.
	//
	MapRemove(ctx context.Context, in *recordpb.MapRemoveRequest) (*recordpb.MapValue, error)

	//
	// Scan all map key-value pairs
	//
	MapRange(ctx context.Context, in *recordpb.MapRangeRequest) (entries <- chan MapEntryEvent, cancel func(), err error)

}


