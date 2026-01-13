package helper

import (
	"golang.org/x/tools/go/packages"
)

var globalPackages []*packages.Package

// SetGlobalPackages is called by runner to provide all loaded packages for cross-package resolution.
func SetGlobalPackages(pkgs []*packages.Package) {
	globalPackages = pkgs
}

// GetGlobalPackages returns all loaded packages.
func GetGlobalPackages() []*packages.Package {
	return globalPackages
}

// FindPackageByPath searches for a package by import path.
func FindPackageByPath(pkgPath string) *packages.Package {
	for _, pkg := range globalPackages {
		if pkg.PkgPath == pkgPath {
			return pkg
		}
	}
	return nil
}
