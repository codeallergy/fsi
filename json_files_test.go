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
	"github.com/codeallergy/fsi"
	"github.com/stretchr/testify/require"
	"github.com/codeallergy/fs"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestJsonWriteAndRead(t *testing.T) {

	fs := fsi.FileService()

	fd, err := ioutil.TempFile(os.TempDir(), "json-test")
	require.NoError(t, err)
	filePath := fd.Name()
	fd.Close()
	os.Remove(filePath)

	// Test Plain
	filePath = filePath + ".json"
	writeJson(t, fs, filePath)
	var buf bytes.Buffer
	writeJsonStream(t, fs.NewJsonStream(&buf, false))

	content, err := ioutil.ReadFile(filePath)
	require.NoError(t, err)
	require.Equal(t, buf.Bytes(), content)

	stream, err := fs.JsonStream(bytes.NewReader(content), false)
	require.NoError(t, err)
	readJsonStream(t, stream)
	readJson(t, fs, filePath)

	os.Remove(filePath)

	// Test GZIP
	filePath = filePath + ".gz"
	writeJson(t, fs, filePath)
	readJson(t, fs, filePath)
	os.Remove(filePath)

}

func writeJson(t *testing.T, fs fs.FileService, filePath string) {

	js, err := fs.NewJsonFile(filePath)
	require.NoError(t, err)

	writeJsonStream(t, js)
}

func writeJsonStream(t *testing.T, js fs.JsonWriter) {

	obj1 := map[string]string {
		"test": "obj1",
	}

	obj2 := map[string]string {
		"test": "obj2",
	}

	err := js.Write(obj1)
	require.NoError(t, err)

	err = js.Write(obj2)
	require.NoError(t, err)

	err = js.Close()
	require.NoError(t, err)
}


func readJson(t *testing.T, fs fs.FileService, filePath string) {

	reader, err := fs.OpenJsonFile(filePath)
	require.NoError(t, err)

	readJsonStream(t, reader)
}

func readJsonStream(t *testing.T, reader fs.JsonReader) {

	obj1 := make(map[string]interface{})

	err := reader.Read(&obj1)
	require.NoError(t, err)

	require.Equal(t, 1, len(obj1))
	require.Equal(t, "obj1", obj1["test"])

	obj2 := make(map[string]interface{})

	err = reader.Read(&obj2)
	require.NoError(t, err)

	require.Equal(t, 1, len(obj2))
	require.Equal(t, "obj2", obj2["test"])

	err = reader.Read(&obj2)
	require.Equal(t, err, io.EOF)

	err = reader.Close()
	require.NoError(t, err)
}

func TestJsonSplit(t *testing.T) {

	fs := fsi.FileService()

	fd, err := ioutil.TempFile(os.TempDir(), "json-test")
	require.NoError(t, err)
	filePath := fd.Name()
	fd.Close()
	os.Remove(filePath)

	jsonFilePath := filePath + ".json"

	jf, err := fs.NewJsonFile(jsonFilePath)
	require.NoError(t, err)

	obj1 := map[string]string {
		"test": "obj1",
	}

	for i := 0; i < 100; i++ {
		err = jf.Write(obj1)
		require.NoError(t, err)
	}

	err = jf.Close()
	require.NoError(t, err)

	parts, err := fs.SplitJsonFile(jsonFilePath, 10, func(i int) string {
		return fmt.Sprintf("%s_part%d.json", filePath, i)
	})
	require.NoError(t, err)

	println(jsonFilePath)
	all, err := ioutil.ReadFile(jsonFilePath)
	require.NoError(t, err)
	//println(string(all))

	err = fs.JoinJsonFiles(jsonFilePath, parts)
	require.NoError(t, err)

	joined, err := ioutil.ReadFile(jsonFilePath)
	require.NoError(t, err)

	require.Equal(t, string(all), string(joined))

	os.Remove(jsonFilePath)
	for _, part := range parts {
		println("RemoveFile: ", part)
		os.Remove(part)
	}
}
