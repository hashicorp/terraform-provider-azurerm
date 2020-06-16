package testsetup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestSetup ensure the vendor contains current repo itself, which might be a dependency of UUT.
// This is not a ideal solution as it cause the vendor directory to be mutated. While it is the
// simplest and efficient solution, allowing user still be able to run `go test`
func TestSetup(m *testing.M) (int, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("getting pwd: %w", err)
	}
	prjdir := pwd
	for i := 0; i < 3; i++ {
		prjdir = filepath.Dir(prjdir)
	}

	vendordir := filepath.Join(prjdir, "vendor")
	oname := filepath.Join(prjdir, "azurerm")
	ndir := filepath.Join(vendordir, "github.com", "terraform-providers", "terraform-provider-azurerm")
	nname := filepath.Join(ndir, "azurerm")

	if err := os.MkdirAll(ndir, os.ModePerm); err != nil {
		return 0, fmt.Errorf("mkdir %s: %w", ndir, err)
	}

	// There is no handy way to test present of a symlink based on
	// error of `os.Symlink` (the wrapped error is syscall.Errno)
	// So we just iterate the file entries to see whether the expected
	// file is present.
	files, err := ioutil.ReadDir(ndir)
	if err != nil {
		return 0, fmt.Errorf("reading dir %s: %w", ndir, err)
	}
	for _, f := range files {
		if f.Name() == "azurerm" {
			return m.Run(), nil
		}
	}

	if err := os.Symlink(oname, nname); err != nil {
		return 0, fmt.Errorf("make symlink %s -> %s: %w", nname, oname, err)
	}

	return m.Run(), nil
}
