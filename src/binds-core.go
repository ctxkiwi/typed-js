// Code generated by go-bindata.
// sources:
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

var _srcCoreTypesTjs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x91\x41\x6e\xb3\x30\x10\x85\xf7\x9c\x62\x96\x3f\x12\x27\x88\xa2\xa0\x7f\xd1\xae\xaa\xa6\x52\x73\x01\x63\x86\x30\x95\x3b\x83\xec\xa1\x22\x8a\x72\xf7\x2a\x36\x21\x82\x84\xae\xd0\xe3\x3d\x7f\xf3\x3c\xce\x6a\x6c\x88\x11\x82\xfa\xde\xea\xf5\x43\x7c\x2c\xc0\x3a\x13\x02\x7c\x46\x05\xe7\x0c\x00\x80\xb8\xc6\x61\xdf\x6c\xa0\xe9\xd9\xfe\x4b\xc1\x1c\x88\x35\xba\xb6\x35\xfe\xbf\x8e\x26\xb1\xe6\x23\x2a\x9a\x0e\xf9\xa8\xed\x66\x0a\x87\xce\x91\x2e\x40\xc6\x7b\x73\xda\x26\xb5\x4b\xa9\xbe\x0a\xea\xef\xc8\xe2\x7a\xbe\x9c\xc0\x97\x6c\xd1\x9d\xfb\xef\x0a\xfd\xad\xfb\x7b\x54\x63\x77\x95\x97\xa1\x13\x46\x56\x32\xee\x4e\x2c\x67\x2d\x55\x5e\x69\xc0\x7a\xdd\x7e\x13\x6b\x1c\xa6\xa5\xcc\xda\x97\x05\x48\xf5\x85\x56\xb7\x86\x4f\xbb\xe5\xb1\x0f\x8f\x96\x02\x09\x3f\x23\x3f\xdc\xa2\x12\x71\x2b\x77\x98\x4d\x5e\x27\xc4\x4d\x3e\x45\xfc\xfd\x84\xcb\x57\xea\xfa\xd0\x8e\xd1\x43\x0e\x3f\x42\x75\xfa\x2d\xdd\x54\xe1\xf0\x38\x3d\x2d\xe2\x36\x7e\x1f\x15\x9c\xb3\x4b\xf6\x1b\x00\x00\xff\xff\x38\xcf\xe1\x3b\x6a\x02\x00\x00")

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

	info := bindataFileInfo{name: "src/core/types.tjs", size: 618, mode: os.FileMode(420), modTime: time.Unix(1610572703, 0)}
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

