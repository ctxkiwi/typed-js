// Code generated by go-bindata.
// sources:
// src/core/globals.tjs
// src/core/types.tjs
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _srcCoreGlobalsTjs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe2\x4a\x49\x4d\xcb\xcc\x4b\x55\xc8\xc9\x4f\x4e\xcc\x51\x28\x2e\x29\x2a\x4d\x2e\x51\x88\x4f\xce\xcf\x2b\xce\xcf\x49\x55\xa8\xe6\x52\x50\x50\x50\xc8\xc9\x4f\xb7\x52\x48\x2b\xcd\x4b\xd6\x28\x2e\x29\xca\xcc\x4b\xd7\x54\x28\xcb\xcf\x4c\xe1\xaa\xe5\x82\x69\x86\xab\x87\xd2\x5c\x80\x00\x00\x00\xff\xff\xda\x28\x21\x24\x56\x00\x00\x00")

func srcCoreGlobalsTjsBytes() ([]byte, error) {
	return bindataRead(
		_srcCoreGlobalsTjs,
		"src/core/globals.tjs",
	)
}

func srcCoreGlobalsTjs() (*asset, error) {
	bytes, err := srcCoreGlobalsTjsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "src/core/globals.tjs", size: 86, mode: os.FileMode(420), modTime: time.Unix(1610822613, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _srcCoreTypesTjs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\x5d\x6a\xf3\x30\x10\x7c\xf7\x29\xf6\xf1\x33\xf8\x04\x21\x04\xf2\x41\xfb\x54\x9a\x42\x73\x01\x59\x5e\xc7\x2a\xea\xae\xd1\x4f\x71\x48\x7a\xf7\x62\xad\x6d\x6a\x3b\xa1\x79\x32\x23\xcd\xcc\xce\xac\x9c\x55\x58\x1b\x42\xf0\xc1\x45\x1d\xfa\x8f\xa1\x53\x01\xda\x2a\xef\xe1\x3d\x21\xb8\x64\x00\x00\x86\x2a\xec\x0e\xf5\x06\xea\x48\xfa\x9f\x10\x73\xa0\xf8\x59\xa2\x4b\x04\xdd\x28\xb7\x0f\xc3\xbd\x9c\xe7\x83\x61\xba\xb7\x48\xa7\xd0\x6c\x7e\x4b\x7c\x6b\x4d\x58\x38\x2a\xe7\xd4\x79\x2b\x68\x27\xac\x58\xfa\xe0\x66\xc6\xc5\xe0\x72\xa5\x68\xed\x35\x92\xb4\xa8\xa6\x79\xdf\xd9\xa2\xd8\x28\x93\x62\xaf\x09\x0d\xc5\x02\x3f\x75\x2d\x13\x52\x30\xca\xce\xa6\xdc\x30\x16\xc1\xb3\xe9\xb0\x7a\x8c\xfa\xc2\x5a\x59\x94\x4d\xce\x9a\x2e\xa2\x17\xc0\xe5\x07\xea\xb0\x55\x74\xde\xdd\x75\x7b\x73\xa8\x8d\x37\x4c\x7f\x0d\x5f\x2d\xa0\x64\xb6\x63\xfd\xff\xcc\x76\x2a\x3f\x8b\x76\x5f\x9f\x9e\x65\x34\xd8\xf7\xe0\xd1\xff\xe2\xc6\xbb\xb7\xd1\x37\x83\xe0\x98\xc3\x17\x9b\x4a\x8e\xb9\x9d\x72\x1c\xd7\x11\x64\x41\x63\x86\x43\x42\x70\x59\xf3\x7a\x87\x74\xfe\x13\x00\x00\xff\xff\x60\xcf\x5a\x68\xdd\x02\x00\x00")

func srcCoreTypesTjsBytes() ([]byte, error) {
	return bindataRead(
		_srcCoreTypesTjs,
		"src/core/types.tjs",
	)
}

func srcCoreTypesTjs() (*asset, error) {
	bytes, err := srcCoreTypesTjsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "src/core/types.tjs", size: 733, mode: os.FileMode(420), modTime: time.Unix(1610817174, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"src/core/globals.tjs": srcCoreGlobalsTjs,
	"src/core/types.tjs": srcCoreTypesTjs,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"src": &bintree{nil, map[string]*bintree{
		"core": &bintree{nil, map[string]*bintree{
			"globals.tjs": &bintree{srcCoreGlobalsTjs, map[string]*bintree{}},
			"types.tjs": &bintree{srcCoreTypesTjs, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

