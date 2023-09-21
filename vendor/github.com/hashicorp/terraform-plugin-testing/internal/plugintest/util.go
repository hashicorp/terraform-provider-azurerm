// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugintest

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func symlinkFile(src string, dest string) error {
	err := os.Symlink(src, dest)

	if err != nil {
		return fmt.Errorf("unable to symlink %q to %q: %w", src, dest, err)
	}

	srcInfo, err := os.Stat(src)

	if err != nil {
		return fmt.Errorf("unable to stat %q: %w", src, err)
	}

	err = os.Chmod(dest, srcInfo.Mode())

	if err != nil {
		return fmt.Errorf("unable to set %q permissions: %w", dest, err)
	}

	return nil
}

// symlinkDirectoriesOnly finds only the first-level child directories in srcDir
// and symlinks them into destDir.
// Unlike symlinkDir, this is done non-recursively in order to limit the number
// of file descriptors used.
func symlinkDirectoriesOnly(srcDir string, destDir string) error {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("unable to stat source directory %q: %w", srcDir, err)
	}

	err = os.MkdirAll(destDir, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("unable to make destination directory %q: %w", destDir, err)
	}

	dirEntries, err := os.ReadDir(srcDir)

	if err != nil {
		return fmt.Errorf("unable to read source directory %q: %w", srcDir, err)
	}

	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}

		srcPath := filepath.Join(srcDir, dirEntry.Name())
		destPath := filepath.Join(destDir, dirEntry.Name())
		err := symlinkFile(srcPath, destPath)

		if err != nil {
			return fmt.Errorf("unable to symlink directory %q to %q: %w", srcPath, destPath, err)
		}
	}

	return nil
}

// CopyFile copies a single file from src to dest.
func CopyFile(src, dest string) error {
	var srcFileInfo os.FileInfo

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("unable to copy: %w", err)
	}

	if srcFileInfo, err = os.Stat(src); err != nil {
		return fmt.Errorf("unable to stat: %w", err)
	}

	return os.Chmod(dest, srcFileInfo.Mode())
}

// CopyDir recursively copies directories and files
// from src to dest.
func CopyDir(src, dest, baseDirName string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("unable to stat: %w", err)
	}

	if err = os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return fmt.Errorf("unable to create dir: %w", err)
	}

	dirEntries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("unable to read dir: %w", err)
	}

	for _, dirEntry := range dirEntries {
		srcFilepath := path.Join(src, dirEntry.Name())
		destFilepath := path.Join(dest, dirEntry.Name())

		if !strings.Contains(srcFilepath, baseDirName) {
			continue
		}

		fi, err := dirEntry.Info()

		if err != nil {
			return fmt.Errorf("unable to get dir entry info: %w", err)
		}

		if dirEntry.IsDir() || fi.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			if err = CopyDir(srcFilepath, destFilepath, baseDirName); err != nil {
				return fmt.Errorf("unable to copy directory: %w", err)
			}
		} else {
			if err = CopyFile(srcFilepath, destFilepath); err != nil {
				return fmt.Errorf("unable to copy file: %w", err)
			}
		}
	}

	return nil
}

// TestExpectTFatal provides a wrapper for logic which should call
// (*testing.T).Fatal() or (*testing.T).Fatalf().
//
// Since we do not want the wrapping test to fail when an expected test error
// occurs, it is required that the testLogic passed in uses
// github.com/mitchellh/go-testing-interface.RuntimeT instead of the real
// *testing.T.
//
// If Fatal() or Fatalf() is not called in the logic, the real (*testing.T).Fatal() will
// be called to fail the test.
func TestExpectTFatal(t *testing.T, testLogic func()) {
	t.Helper()

	var recoverIface interface{}

	func() {
		defer func() {
			recoverIface = recover()
		}()

		testLogic()
	}()

	if recoverIface == nil {
		t.Fatalf("expected t.Fatal(), got none")
	}

	recoverStr, ok := recoverIface.(string)

	if !ok {
		t.Fatalf("expected string from recover(), got: %v (%T)", recoverIface, recoverIface)
	}

	// this string is hardcoded in github.com/mitchellh/go-testing-interface
	if !strings.HasPrefix(recoverStr, "testing.T failed, see logs for output") {
		t.Fatalf("expected t.Fatal(), got: %s", recoverStr)
	}
}
