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

package fsi

import (
	"github.com/codeallergy/fs"
	"google.golang.org/protobuf/encoding/protojson"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// default size is 64kb, possible to overwrite
var DefaultBufferSize = 64 * 1024

type fileServiceImpl struct {
	bufferSize int // read/write block buffer size
	marshaler  runtime.JSONPb
}

func FileService() fs.FileService {
	return &fileServiceImpl{
		bufferSize: DefaultBufferSize,
		marshaler: runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}
}

func (t *fileServiceImpl) BufferSize() int {
	return t.bufferSize
}

func (t *fileServiceImpl) SetBufferSize(size int) {
	t.bufferSize = size
}

func (t *fileServiceImpl) MarshalOptions() protojson.MarshalOptions {
	return t.marshaler.MarshalOptions
}

func (t *fileServiceImpl) SetMarshalOptions(opt protojson.MarshalOptions) {
	t.marshaler.MarshalOptions = opt
}

func (t *fileServiceImpl) UnmarshalOptions() protojson.UnmarshalOptions {
	return t.marshaler.UnmarshalOptions
}

func (t * fileServiceImpl) SetUnmarshalOptions(opt protojson.UnmarshalOptions) {
	t.marshaler.UnmarshalOptions = opt
}



