// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

type (
	BasicMapReader         = schema.BasicMapReader
	MapFieldReader         = schema.MapFieldReader
	MapFieldWriter         = schema.MapFieldWriter
	Resource               = schema.Resource
	ResourceData           = schema.ResourceData
	ResourceDiff           = schema.ResourceDiff
	SchemaDiffSuppressFunc = schema.SchemaDiffSuppressFunc
	StateUpgrader          = schema.StateUpgrader
	SchemaValidateFunc     = func(interface{}, string) ([]string, []error)
	ValueType              = schema.ValueType
)

type (
	StateChangeConf  = retry.StateChangeConf
	StateRefreshFunc = retry.StateRefreshFunc
)

type (
	// lint:ignore SA1019 SDKv2 migration - staticcheck's own linter directives are currently being ignored under golanci-lint
	CreateFunc = schema.CreateFunc //nolint:staticcheck
	DeleteFunc = schema.DeleteFunc //nolint:staticcheck
	ExistsFunc = schema.ExistsFunc //nolint:staticcheck
	ReadFunc   = schema.ReadFunc   //nolint:staticcheck
	UpdateFunc = schema.UpdateFunc //nolint:staticcheck
)
