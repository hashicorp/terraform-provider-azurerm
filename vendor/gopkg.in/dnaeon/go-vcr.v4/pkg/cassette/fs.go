// Copyright (c) 2015-2024 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer
//    in this position and unchanged.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package cassette

import (
	"os"
	"path/filepath"
)

// FS defines a generic filesystem interface.
// It allows to redefine storage without depending on a specific filesystem implementation.
type FS interface {
	// ReadFile reads the contents of the file named by 'name'.
	ReadFile(name string) ([]byte, error)

	// WriteFile writes 'data' to the file named by 'name'.
	WriteFile(name string, data []byte) error

	// IsFileExists checks whether a file with the given 'name' exists.
	IsFileExists(name string) bool
}

// NewDiskFS creates and returns a new FS implementation backed by the local disk filesystem.
func NewDiskFS() FS {
	return &diskFS{}
}

type diskFS struct{}

func (fs *diskFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs *diskFS) WriteFile(name string, data []byte) error {
	cassetteDir := filepath.Dir(name)
	if _, err := os.Stat(cassetteDir); os.IsNotExist(err) {
		if err = os.MkdirAll(cassetteDir, 0o755); err != nil {
			return err
		}
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.Write(data)
	return err
}

func (fs *diskFS) IsFileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
