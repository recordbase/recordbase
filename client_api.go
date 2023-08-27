/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package recordbase

import (
	"context"
	"github.com/codeallergy/glue"
	"github.com/recordbase/recordbasepb"
	"reflect"
)

const (
	ServiceName = "recordbase.RecordService"
)

var ClientClass = reflect.TypeOf((*Client)(nil)).Elem()

type RecordEntryEvent struct {
	Entry      *recordbasepb.RecordEntry
	Err        error
}

type MapEntryEvent struct {
	Entry      *recordbasepb.MapEntry
	Err        error
}

type FileContentEvent struct {
	Content    *recordbasepb.FileContent
	Err        error
}

type Client interface {
	glue.DisposableBean

	//
	// Gets metadata about using attributes
	//
	GetInfo(ctx context.Context, in *recordbasepb.TenantRequest) (*recordbasepb.Info, error)

	//
	// Quick user lookup request
	//
	Lookup(ctx context.Context, in *recordbasepb.LookupRequest) (*recordbasepb.RecordEntry, error)

	//
	// Search users by indexed non-unique attributes
	//
	Search(ctx context.Context, in *recordbasepb.SearchRequest) (entries <- chan RecordEntryEvent, cancel func(), err error)

	//
	// Get user with all attributes
	//
	Get(ctx context.Context, in *recordbasepb.GetRequest) (*recordbasepb.RecordEntry, error)

	//
	// Create user, returns new user_id
	//
	Create(ctx context.Context, in *recordbasepb.CreateRequest) (*recordbasepb.CreateResponse, error)

	//
	// Delete user request (sets TTL to all PII data for particular user)
	//
	Delete(ctx context.Context, in *recordbasepb.DeleteRequest) error

	//
	// Update user attributes
	//
	Update(ctx context.Context, in *recordbasepb.UpdateRequest) error

	//
	// Upload File
	//
	UploadFile(c context.Context) (sink chan <- *recordbasepb.UploadFileContent, _ <- chan error)

	//
	// Download File
	//
	DownloadFile(ctx context.Context, in *recordbasepb.DownloadFileRequest) (entries <- chan FileContentEvent, cancel func(), err error)

	//
	// Delete File
	//
	DeleteFile(ctx context.Context, in *recordbasepb.DeleteFileRequest) error

	//
	// Scan users
	//
	Scan(ctx context.Context, in *recordbasepb.ScanRequest) (entries <- chan RecordEntryEvent, cancel func(), err error)

	//
	// Allocate user id range
	//
	AddKeyRange(ctx context.Context, in *recordbasepb.KeyRange) error

	//
	// Gets user id ranges and etc
	//
	GetKeyCapacity(ctx context.Context, in *recordbasepb.TenantRequest) (*recordbasepb.KeyCapacity, error)

	//
	// Get map value associated with the record
	//
	MapGet(ctx context.Context, in *recordbasepb.MapGetRequest) (*recordbasepb.MapEntry, error)

	//
	// Put map value associated with the record. Returns old value.
	//
	MapPut(ctx context.Context, in *recordbasepb.MapPutRequest) error

	//
	// Remove map value associated with the record. Returns old value.
	//
	MapRemove(ctx context.Context, in *recordbasepb.MapRemoveRequest) error

	//
	// Scan all map key-value pairs
	//
	MapRange(ctx context.Context, in *recordbasepb.MapRangeRequest) (entries <- chan MapEntryEvent, cancel func(), err error)

	//
	// Get bin value from the record
	//
	BinGet(ctx context.Context, in *recordbasepb.BinGetRequest) (*recordbasepb.BinEntry, error)

	//
	// Put bin value to the record
	//
	BinPut(ctx context.Context, in *recordbasepb.BinPutRequest) error

	//
	// Remove bin value from the record
	//
	BinRemove(ctx context.Context, in *recordbasepb.BinRemoveRequest) error


}


