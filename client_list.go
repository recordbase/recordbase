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
	"errors"
	"github.com/recordbase/recordpb"
	"sync"
)

var ErrInstanceNotExist = errors.New("instance not exist")

type ClientList struct {
	mu sync.RWMutex
	list  []Client
}

func (t *ClientList) Add(client Client) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	pos := len(t.list)
	t.list = append(t.list, client)
	return pos
}

func (t *ClientList) Get(pos int) Client {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if pos >= 0 && pos < len(t.list) {
		return t.list[pos]
	}
	return dummyClientImpl{}
}

func (t *ClientList) Remove(pos int) Client {
	t.mu.Lock()
	defer t.mu.Unlock()
	if pos >= 0 && pos < len(t.list) {
		cli := t.list[pos]
		t.list[pos] = dummyClientImpl{}
		return cli
	}
	return dummyClientImpl{}
}

type dummyClientImpl struct {
}

func (d dummyClientImpl) Destroy() error {
	return ErrInstanceNotExist
}

func (d dummyClientImpl) GetCounts(ctx context.Context, in *recordpb.TenantRequest) (*recordpb.Counts, error) {
	return nil, ErrInstanceNotExist
}

func (d dummyClientImpl) Lookup(ctx context.Context, in *recordpb.LookupRequest) (*recordpb.RecordEntry, error) {
	return nil, ErrInstanceNotExist
}

func (d dummyClientImpl) Search(ctx context.Context, in *recordpb.SearchRequest) (entries <-chan RecordEntryEvent, cancel func(), err error) {
	return nil, nil, ErrInstanceNotExist
}

func (d dummyClientImpl) Get(ctx context.Context, in *recordpb.GetRequest) (*recordpb.RecordEntry, error) {
	return nil, ErrInstanceNotExist
}

func (d dummyClientImpl) Create(ctx context.Context, in *recordpb.CreateRequest) (*recordpb.CreateResponse, error) {
	return nil, ErrInstanceNotExist
}

func (d dummyClientImpl) Delete(ctx context.Context, in *recordpb.DeleteRequest) error {
	return ErrInstanceNotExist
}

func (d dummyClientImpl) Update(ctx context.Context, in *recordpb.UpdateRequest) error {
	return ErrInstanceNotExist
}

func (d dummyClientImpl) Scan(ctx context.Context, in *recordpb.ScanRequest) (entries <-chan RecordEntryEvent, cancel func(), err error) {
	return nil, nil, ErrInstanceNotExist
}

func (d dummyClientImpl) AddKeyRange(ctx context.Context, in *recordpb.KeyRange) error {
	return ErrInstanceNotExist
}

func (d dummyClientImpl) GetKeyCapacity(ctx context.Context, in *recordpb.TenantRequest) (*recordpb.KeyCapacity, error) {
	return nil, ErrInstanceNotExist
}

func (d dummyClientImpl) MapGet(ctx context.Context, in *recordpb.MapGetRequest) (*recordpb.MapEntry, error) {
	return nil, ErrInstanceNotExist
}

func (d dummyClientImpl) MapPut(ctx context.Context, in *recordpb.MapPutRequest) (*recordpb.MapValue, error) {
	return nil, ErrInstanceNotExist
}

func (d dummyClientImpl) MapRemove(ctx context.Context, in *recordpb.MapRemoveRequest) (*recordpb.MapValue, error) {
	return nil, ErrInstanceNotExist
}

func (d dummyClientImpl) MapRange(ctx context.Context, in *recordpb.MapRangeRequest) (entries <-chan MapEntryEvent, cancel func(), err error) {
	return nil, nil, ErrInstanceNotExist
}


