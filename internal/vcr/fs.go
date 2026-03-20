// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vcr

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// gzipFS is a cassette.FS implementation that reads and writes gzip-compressed files.
type gzipFS struct{}

func (fs *gzipFS) ReadFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

func (fs *gzipFS) WriteFile(name string, data []byte) error {
	dir := filepath.Dir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())
	return err
}

func (fs *gzipFS) IsFileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
