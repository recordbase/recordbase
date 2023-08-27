/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package recordbase

import (
	"context"
	"github.com/recordbase/recordbasepb"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"io"
	"sync"
)

type implClient struct {

	conn             *grpc.ClientConn
	client           recordbasepb.RecordServiceClient

	nextHandle       atomic.Int64
	cancelFunctions  sync.Map

	closeOnce    sync.Once

}

func Create(conn *grpc.ClientConn) Client {
	return &implClient{
		conn: conn,
		client: recordbasepb.NewRecordServiceClient(conn),
	}
}

func (t *implClient) GetInfo(ctx context.Context, in *recordbasepb.TenantRequest) (*recordbasepb.Info, error) {

	ctx, cancel := context.WithCancel(ctx)

	handle := t.addCancelFn(cancel)
	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	return t.client.GetInfo(ctx, in)
}

func (t *implClient) Lookup(ctx context.Context, in *recordbasepb.LookupRequest) (*recordbasepb.RecordEntry, error) {

	ctx, cancel := context.WithCancel(ctx)

	handle := t.addCancelFn(cancel)
	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	return t.client.Lookup(ctx, in)
}

func (t *implClient) Search(c context.Context, in *recordbasepb.SearchRequest) (entries <-chan RecordEntryEvent, cancel func(), err error) {

	ctx, cancel := context.WithCancel(c)
	handle := t.addCancelFn(cancel)

	stream, err := t.client.Search(ctx, in)

	if err != nil {
		t.removeCancelFn(handle)
		return nil, nil, err
	}

	resultCh := make(chan RecordEntryEvent)

	go func() {

		defer func() {
			t.removeCancelFn(handle)
			close(resultCh)
		}()

		for ctx.Err() == nil {
			entry, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					resultCh <- RecordEntryEvent { Err: err }
				}
				break
			}
			resultCh <- RecordEntryEvent { Entry: entry }
		}
	}()

	return resultCh, func() {
		t.removeCancelFn(handle)
		cancel()
	},nil

}

func (t *implClient) Get(ctx context.Context, in *recordbasepb.GetRequest) (*recordbasepb.RecordEntry, error) {

	ctx, cancel := context.WithCancel(ctx)

	handle := t.addCancelFn(cancel)
	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	return t.client.Get(ctx, in)
}

func (t *implClient) Create(ctx context.Context, in *recordbasepb.CreateRequest) (*recordbasepb.CreateResponse, error) {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	return t.client.Create(ctx, in)
}

func (t *implClient) Delete(ctx context.Context, in *recordbasepb.DeleteRequest) error {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	_, err := t.client.Delete(ctx, in)
	return  err
}

func (t *implClient) Update(ctx context.Context, in *recordbasepb.UpdateRequest) error {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	_, err := t.client.Update(ctx, in)
	return  err
}

func (t *implClient) UploadFile(c context.Context) (chan <- *recordbasepb.UploadFileContent, <- chan error) {

	errCh := make(chan error)

	ctx, cancel := context.WithCancel(c)
	handle := t.addCancelFn(cancel)

	stream, errCall := t.client.UploadFile(ctx)
	if errCall != nil {
		t.removeCancelFn(handle)
		errCh <- errCall
		return nil, errCh
	}

	sink := make(chan *recordbasepb.UploadFileContent)

	go func() {

		for content := range sink {
			err := stream.Send(content)
			if err != nil {
				errCh <- err
				return
			}
		}

		err := stream.CloseSend()
		if err != nil {
			errCh <- err
		}

	}()

	return sink, errCh

}

func (t *implClient) DownloadFile(c context.Context, in *recordbasepb.DownloadFileRequest)  (entries <- chan FileContentEvent, cancel func(), err error) {

	ctx, cancel := context.WithCancel(c)
	handle := t.addCancelFn(cancel)

	stream, err := t.client.DownloadFile(ctx, in)

	if err != nil {
		t.removeCancelFn(handle)
		return nil, nil, err
	}

	resultCh := make(chan FileContentEvent)

	go func() {

		defer func() {
			t.removeCancelFn(handle)
			close(resultCh)
		}()

		for ctx.Err() == nil {
			entry, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					resultCh <- FileContentEvent { Err: err }
				}
				break
			}
			resultCh <- FileContentEvent { Content: entry }
		}
	}()

	return resultCh, func() {
		t.removeCancelFn(handle)
		cancel()
	},nil

}

func (t *implClient) DeleteFile(ctx context.Context, in *recordbasepb.DeleteFileRequest) error {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	_, err := t.client.DeleteFile(ctx, in)
	return  err
}

func (t *implClient) Scan(c context.Context, in *recordbasepb.ScanRequest) (entries <-chan RecordEntryEvent, cancel func(), err error) {

	ctx, cancel := context.WithCancel(c)
	handle := t.addCancelFn(cancel)

	stream, err := t.client.Scan(ctx, in)

	if err != nil {
		t.removeCancelFn(handle)
		return nil, nil, err
	}

	resultCh := make(chan RecordEntryEvent)

	go func() {

		defer func() {
			t.removeCancelFn(handle)
			close(resultCh)
		}()

		for ctx.Err() == nil {
			entry, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					resultCh <- RecordEntryEvent { Err: err }
				}
				break
			}
			resultCh <- RecordEntryEvent { Entry: entry }
		}
	}()

	return resultCh, func() {
		t.removeCancelFn(handle)
		cancel()
	},nil
}

func (t *implClient) AddKeyRange(ctx context.Context, in *recordbasepb.KeyRange) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	_, err = t.client.AddKeyRange(ctx, in)
	return err

}

func (t *implClient) GetKeyCapacity(ctx context.Context, in *recordbasepb.TenantRequest) (*recordbasepb.KeyCapacity, error) {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	return t.client.GetKeyCapacity(ctx, in)
}

func (t *implClient) MapGet(ctx context.Context, in *recordbasepb.MapGetRequest) (*recordbasepb.MapEntry, error) {

	ctx, cancel := context.WithCancel(ctx)

	handle := t.addCancelFn(cancel)
	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	return t.client.MapGet(ctx, in)

}

func (t *implClient) MapPut(ctx context.Context, in *recordbasepb.MapPutRequest) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	_, err = t.client.MapPut(ctx, in)
	return err

}

func (t *implClient) MapRemove(ctx context.Context, in *recordbasepb.MapRemoveRequest) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	_, err = t.client.MapRemove(ctx, in)
	return err

}

func (t *implClient) MapRange(c context.Context, in *recordbasepb.MapRangeRequest) (entries <-chan MapEntryEvent, cancel func(), err error) {

	ctx, cancel := context.WithCancel(c)
	handle := t.addCancelFn(cancel)

	stream, err := t.client.MapRange(ctx, in)

	if err != nil {
		t.removeCancelFn(handle)
		return nil, nil, err
	}

	resultCh := make(chan MapEntryEvent)

	go func() {

		defer func() {
			t.removeCancelFn(handle)
			close(resultCh)
		}()

		for ctx.Err() == nil {
			entry, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					resultCh <- MapEntryEvent { Err: err }
				}
				break
			}
			resultCh <- MapEntryEvent { Entry: entry }
		}
	}()

	return resultCh, func() {
		t.removeCancelFn(handle)
		cancel()
	},nil

}

func (t *implClient) BinGet(ctx context.Context, in *recordbasepb.BinGetRequest) (*recordbasepb.BinEntry, error) {

	ctx, cancel := context.WithCancel(ctx)

	handle := t.addCancelFn(cancel)
	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	return t.client.BinGet(ctx, in)

}

func (t *implClient) BinPut(ctx context.Context, in *recordbasepb.BinPutRequest) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	_, err = t.client.BinPut(ctx, in)
	return err

}

func (t *implClient) BinRemove(ctx context.Context, in *recordbasepb.BinRemoveRequest) (err error) {

	ctx, cancel := context.WithCancel(ctx)
	handle := t.addCancelFn(cancel)

	defer func() {
		t.removeCancelFn(handle)
		cancel()
	}()

	_, err = t.client.BinRemove(ctx, in)
	return err

}

func (t *implClient) addCancelFn(cancelFn func()) int64 {
	handle := t.nextHandle.Inc()
	t.cancelFunctions.Store(handle, cancelFn)
	return handle
}

func (t *implClient) removeCancelFn(handle int64) {
	t.cancelFunctions.Delete(handle)
}

func (t *implClient) Destroy() (err error) {
	t.closeOnce.Do(func() {

		err = t.conn.Close()

		t.cancelFunctions.Range(func(key, value interface{}) bool {
			if closeFn, ok := value.(func()); ok {
				closeFn()
			}
			return true
		})

	})
	return nil
}

