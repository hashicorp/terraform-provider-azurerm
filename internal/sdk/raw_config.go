// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
)

// GetRawConfig is a helper to retrieve the RawConfig from a ResourceMetaData object
//
// Note:
// This method is experimental and not meant for general use.
// Pull requests using this method, by authors not part of the HashiCorp AzureRM Provider team, will be declined at this time.
func (rmd ResourceMetaData) GetRawConfig() (cty.Value, error) {
	switch {
	case rmd.ResourceData != nil:
		return rmd.ResourceData.GetRawConfig(), nil
	case rmd.ResourceDiff != nil:
		return rmd.ResourceDiff.GetRawConfig(), nil
	}

	// It *shouldn't* be possible to reach this return statement
	return cty.NilVal, errors.New("internal error: both `ResourceData` and `ResourceDiff` were nil, unable to retrieve RawConfig")
}

// GetRawConfigAt is a helper to retrieve the RawConfig value for a specific argument path from a ResourceMetaData object
//
// Note:
// This method is experimental and not meant for general use.
// Pull requests using this method, by authors not part of the HashiCorp AzureRM Provider team, will be declined at this time.
func (rmd ResourceMetaData) GetRawConfigAt(key string) (cty.Value, error) {
	ctyPath := ConstructCtyPath(key)

	msg := "retrieving value at path `%s`: %+v"

	switch {
	case rmd.ResourceData != nil:
		val, diags := rmd.ResourceData.GetRawConfigAt(ctyPath)
		if diags.HasError() {
			return cty.NilVal, fmt.Errorf(msg, key, diags)
		}
		return val, nil
	case rmd.ResourceDiff != nil:
		val, diags := rmd.ResourceDiff.GetRawConfigAt(ctyPath)
		if diags.HasError() {
			return cty.NilVal, fmt.Errorf(msg, key, diags)
		}
		return val, nil
	}

	// It *shouldn't* be possible to reach this return statement
	return cty.NilVal, errors.New("internal error: both `ResourceData` and `ResourceDiff` were nil, unable to retrieve RawConfig")
}

// GetRawConfigAsValueMap is a helper to retrieve the RawConfig map handling all safety checks
//
// Note:
// This method is experimental and not meant for general use.
// Pull requests using this method, by authors not part of the HashiCorp AzureRM Provider team, will be declined at this time.
func (rmd ResourceMetaData) GetRawConfigAsValueMap() (map[string]cty.Value, error) {
	rawConfig, err := rmd.GetRawConfig()
	if err != nil {
		return nil, err
	}

	return asValueMap(rawConfig)
}

// IsKnownAt returns whether a value at a specified path is known
// If there were any errors while retrieving the value, this method will return `false`
//
// Note:
// This method is experimental and not meant for general use.
// Pull requests using this method, by authors not part of the HashiCorp AzureRM Provider team, will be declined at this time.
func (rmd ResourceMetaData) IsKnownAt(key string) bool {
	val, err := rmd.GetRawConfigAt(key)
	if err != nil {
		return false
	}

	return val.IsKnown()
}

// IsWhollyKnownAt returns whether a value at a specified path and its nested attributes are known
// If there were any errors while retrieving the value, this method will return `false`
// This should be used for properties of `TypeSet`, `TypeList`, or `TypeMap`.
//
// Note:
// This method is experimental and not meant for general use.
// Pull requests using this method, by authors not part of the HashiCorp AzureRM Provider team, will be declined at this time.
func (rmd ResourceMetaData) IsWhollyKnownAt(key string) bool {
	val, err := rmd.GetRawConfigAt(key)
	if err != nil {
		return false
	}

	return val.IsWhollyKnown()
}

// IsNullAt returns whether a value at a specified path is known
// If there were any errors while retrieving the value, this method will return `false`
//
// Note:
// This method is experimental and not meant for general use.
// Pull requests using this method, by authors not part of the HashiCorp AzureRM Provider team, will be declined at this time.
func (rmd ResourceMetaData) IsNullAt(key string) bool {
	val, err := rmd.GetRawConfigAt(key)
	if err != nil {
		return false
	}

	return val.IsNull()
}

// asValueMap is a convenience function that handles all the safety checks to prevent panics before calling `[cty.Value].AsValueMap()`
//
// Note:
// This function is experimental and not meant for general use.
// Pull requests using this function, by authors not part of the HashiCorp AzureRM Provider team, will be declined at this time.
func asValueMap(val cty.Value) (map[string]cty.Value, error) {
	if val.IsNull() {
		return nil, errors.New("internal error: provided cty.Value was null")
	}

	if !val.IsKnown() {
		return nil, errors.New("internal error: provided cty.Value was unknown")
	}

	if !val.CanIterateElements() {
		return nil, errors.New("internal error: provided cty.Value does not support the ElementIterator method")
	}

	return val.AsValueMap(), nil
}

// ConstructCtyPath takes a string and converts it to a `cty.Path` for use with `GetRawConfigAt`
// e.g. `identity.0.type`
//
// Note:
// This function is experimental and not meant for general use.
// Pull requests using this function, by authors not part of the HashiCorp AzureRM Provider team, will be declined at this time.
func ConstructCtyPath(key string) cty.Path {
	p := cty.Path{}

	for _, segment := range strings.Split(key, ".") {
		if n, err := strconv.Atoi(segment); err == nil {
			p = p.IndexInt(n)
			continue
		}
		p = p.GetAttr(segment)
	}

	return p
}
