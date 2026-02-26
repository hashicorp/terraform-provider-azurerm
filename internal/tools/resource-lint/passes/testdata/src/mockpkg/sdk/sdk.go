// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"

	"testdata/src/mockpkg/pluginsdk"
)

// Mock types for azurerm SDK

type ResourceMetaData struct {
	ResourceData *pluginsdk.ResourceData
}

func (r *ResourceMetaData) Decode(v interface{}) error {
	return nil
}

type ResourceFunc struct {
	Func    func(context.Context, ResourceMetaData) error
	Timeout int
}

type Resource interface {
	ResourceType() string
	ModelObject() interface{}
	Arguments() map[string]*pluginsdk.Schema
	Attributes() map[string]*pluginsdk.Schema
	Create() ResourceFunc
	Read() ResourceFunc
	Delete() ResourceFunc
}

type ResourceWithUpdate interface {
	Resource
	Update() ResourceFunc
}
