// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package autocert

import (
	"context"
	"errors"
	"os"
	"path/filepath"
)

// ErrCacheMiss is returned when a certificate is not found in cache.
var ErrCacheMiss = errors.New("acme/autocert: certificate cache miss")

// Cache is used by Manager to store and retrieve previously obtained certificates
// and other account data as opaque blobs.
//
// Cache implementations should not rely on the certs naming pattern. Keys can
// include any printable ASCII characters, except the following: \/:*?"<>|
type Cache interface {
	// Get returns a certificate data for the specified certs.
	// If there's no such certs, Get returns ErrCacheMiss.
	Get(ctx context.Context, key string) ([]byte, error)

	// Put stores the data in the cache under the specified certs.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, Get, results in the original data.
	Put(ctx context.Context, key string, data []byte) error

	// Delete removes a certificate data from the cache under the specified certs.
	// If there's no such certs in the cache, Delete returns nil.
	Delete(ctx context.Context, key string) error
}

// DirCache implements Cache using a directory on the local filesystem.
// If the directory does not exist, it will be created with 0700 permissions.
type DirCache string

// Get reads a certificate data from the specified file name.
func (d DirCache) Get(ctx context.Context, name string) ([]byte, error) {
	name = filepath.Join(string(d), filepath.Clean("/"+name))
	var (
		data []byte
		err  error
		done = make(chan struct{})
	)
	go func() {
		data, err = os.ReadFile(name)
		close(done)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	}
	if os.IsNotExist(err) {
		return nil, ErrCacheMiss
	}
	return data, err
}

// Put writes the certificate data to the specified file name.
// The file will be created with 0600 permissions.
func (d DirCache) Put(ctx context.Context, name string, data []byte) error {
	if err := os.MkdirAll(string(d), 0700); err != nil {
		return err
	}

	done := make(chan struct{})
	var err error
	go func() {
		defer close(done)
		var tmp string
		if tmp, err = d.writeTempFile(name, data); err != nil {
			return
		}
		defer os.Remove(tmp)
		select {
		case <-ctx.Done():
			// Don't overwrite the file if the context was canceled.
		default:
			newName := filepath.Join(string(d), filepath.Clean("/"+name))
			err = os.Rename(tmp, newName)
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}
	return err
}

// Delete removes the specified file name.
func (d DirCache) Delete(ctx context.Context, name string) error {
	name = filepath.Join(string(d), filepath.Clean("/"+name))
	var (
		err  error
		done = make(chan struct{})
	)
	go func() {
		err = os.Remove(name)
		close(done)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// writeTempFile writes b to a temporary file, closes the file and returns its path.
func (d DirCache) writeTempFile(prefix string, b []byte) (name string, reterr error) {
	// TempFile uses 0600 permissions
	f, err := os.CreateTemp(string(d), prefix)
	if err != nil {
		return "", err
	}
	defer func() {
		if reterr != nil {
			os.Remove(f.Name())
		}
	}()
	if _, err := f.Write(b); err != nil {
		f.Close()
		return "", err
	}
	return f.Name(), f.Close()
}
