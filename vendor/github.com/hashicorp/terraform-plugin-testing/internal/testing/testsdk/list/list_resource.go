// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"context"
	"iter"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type ListResource interface {
	Schema(context.Context, SchemaRequest, *SchemaResponse)
	List(context.Context, ListRequest, *ListResultsStream)
	ValidateListConfig(context.Context, ValidateListConfigRequest, *ValidateListConfigResponse)
}

type ListRequest struct {
	TypeName string
	// Config is the configuration the user supplied for listing resource
	// instances.
	Config tftypes.Value

	// IncludeResource indicates whether the provider should populate the
	// [ListResult.Resource] field.
	IncludeResource bool

	// Limit specifies the maximum number of results that Terraform is
	// expecting.
	Limit int64

	ResourceSchema         *tfprotov6.Schema
	ResourceIdentitySchema *tfprotov6.ResourceIdentitySchema
}

type ListResultsStream struct {
	Results iter.Seq[ListResult]
}

type ListResult struct {
	DisplayName string
	Identity    *tftypes.Value
	Resource    *tftypes.Value
	Diagnostics []*tfprotov6.Diagnostic
}

type ValidateListConfigRequest struct {
	Config tftypes.Value
}

type ValidateListConfigResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
}

type SchemaRequest struct{}

type SchemaResponse struct {
	Diagnostics []*tfprotov6.Diagnostic
	Schema      *tfprotov6.Schema
}
