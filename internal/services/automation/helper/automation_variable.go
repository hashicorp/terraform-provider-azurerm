// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BooleanVariable struct {
	ID          string `tfschema:"id"`
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	IsEncrypted bool   `tfschema:"encrypted"`
	Value       bool   `tfschema:"value"`
}

type DateTimeVariable struct {
	ID          string `tfschema:"id"`
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	IsEncrypted bool   `tfschema:"encrypted"`
	Value       string `tfschema:"value"`
}

type EncryptedVariable struct {
	ID          string `tfschema:"id"`
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	IsEncrypted bool   `tfschema:"encrypted"`
}

type IntegerVariable struct {
	ID          string `tfschema:"id"`
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	IsEncrypted bool   `tfschema:"encrypted"`
	Value       int64  `tfschema:"value"`
}

type NullVariable struct {
	ID          string `tfschema:"id"`
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	IsEncrypted bool   `tfschema:"encrypted"`
}

type ObjectVariable struct {
	ID          string `tfschema:"id"`
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	IsEncrypted bool   `tfschema:"encrypted"`
	Value       string `tfschema:"value"`
}

type StringVariable struct {
	ID          string `tfschema:"id"`
	Name        string `tfschema:"name"`
	Description string `tfschema:"description"`
	IsEncrypted bool   `tfschema:"encrypted"`
	Value       string `tfschema:"value"`
}

func DataSourceAutomationVariableCommonSchema(attType pluginsdk.ValueType) map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"encrypted": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"value": {
			Type:     attType,
			Computed: true,
		},
	}
}
