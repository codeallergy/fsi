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
package fsi_test

import (
	"bytes"
	"fmt"
	"github.com/codeallergy/fs"
	"github.com/codeallergy/fsi"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestProtoWriteAndRead(t *testing.T) {

	fs := fsi.FileService()

	fd, err := ioutil.TempFile(os.TempDir(), "proto-test")
	require.NoError(t, err)
	filePath := fd.Name()
	fd.Close()
	os.Remove(filePath)

	// Test Plain
	filePath = filePath + ".pb"
	writeProto(t, fs, filePath)
	var buf bytes.Buffer
	writeProtoStream(t, fs.NewProtoStream(&buf, false))

	content, err := ioutil.ReadFile(filePath)
	require.NoError(t, err)
	require.Equal(t, buf.Bytes(), content)

	stream, err := fs.ProtoStream(bytes.NewReader(content), false)
	require.NoError(t, err)
	readProtoStream(t, stream)
	readProto(t, fs, filePath)

	os.Remove(filePath)

	// Test GZIP
	filePath = filePath + ".gz"
	writeProto(t, fs, filePath)
	readProto(t, fs, filePath)
	os.Remove(filePath)

}

func writeProto(t *testing.T, fs fs.FileService, filePath string) {

	pf, err := fs.NewProtoFile(filePath)
	require.NoError(t, err)

	writeProtoStream(t, pf)
}

func writeProtoStream(t *testing.T, pf fs.ProtoWriter) {

	obj1 := &Domain{
		Domain:                 "obj1",
	}

	obj2 := &Domain{
		Domain:                 "obj2",
	}

	_, err := pf.Write(obj1)
	require.NoError(t, err)

	_, err = pf.Write(obj2)
	require.NoError(t, err)

	err = pf.Close()
	require.NoError(t, err)
}

func readProto(t *testing.T, fs fs.FileService, filePath string) {

	reader, err := fs.OpenProtoFile(filePath)
	require.NoError(t, err)

	readProtoStream(t, reader)

}

func readProtoStream(t *testing.T, reader fs.ProtoReader) {

	var obj1 Domain

	err := reader.ReadTo(&obj1)
	require.NoError(t, err)

	require.Equal(t, "obj1", obj1.Domain)

	var obj2 Domain

	err = reader.ReadTo(&obj2)
	require.NoError(t, err)

	require.Equal(t, "obj2", obj2.Domain)

	err = reader.ReadTo(&obj2)
	require.Equal(t, err, io.EOF)
}


func TestProtoSplit(t *testing.T) {

	fs := fsi.FileService()

	fd, err := ioutil.TempFile(os.TempDir(), "proto-test")
	require.NoError(t, err)
	filePath := fd.Name()
	fd.Close()
	os.Remove(filePath)

	protoFilePath := filePath + ".pb"

	pf, err := fs.NewProtoFile(protoFilePath)
	require.NoError(t, err)

	obj1 := &Domain{
		Domain:                 "obj1",
	}

	for i := 0; i < 100; i++ {
		_, err = pf.Write(obj1)
		require.NoError(t, err)
	}

	err = pf.Close()
	require.NoError(t, err)

	parts, err := fs.SplitProtoFile(protoFilePath, obj1, 10, func(i int) string {
		return fmt.Sprintf("%s_part%d.pb", filePath, i)
	})
	require.NoError(t, err)

	println(protoFilePath)
	all, err := ioutil.ReadFile(protoFilePath)
	require.NoError(t, err)
	//println(string(all))

	err = fs.JoinProtoFiles(protoFilePath, obj1, parts)
	require.NoError(t, err)

	joined, err := ioutil.ReadFile(protoFilePath)
	require.NoError(t, err)

	require.Equal(t, all, joined)

	os.Remove(protoFilePath)
	for _, part := range parts {
		println("RemoveFile: ", part)
		os.Remove(part)
	}
}
