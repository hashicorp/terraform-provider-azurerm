// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

var _ resource.Resource = Resource{}

type Resource struct {
	CreateResponse *resource.CreateResponse
	// Some tests need more control over logic run during Create, the struct defined CreateResponse
	// will be passed to this function if available
	CreateFunc func(context.Context, resource.CreateRequest, *resource.CreateResponse)

	DeleteResponse      *resource.DeleteResponse
	ImportStateResponse *resource.ImportStateResponse

	// Planning happens multiple ways during a single TestStep, so statically
	// defining only the response is very problematic.
	PlanChangeFunc func(context.Context, resource.PlanChangeRequest, *resource.PlanChangeResponse)

	ReadResponse           *resource.ReadResponse
	IdentitySchemaResponse *resource.IdentitySchemaResponse
	SchemaResponse         *resource.SchemaResponse
	UpdateResponse         *resource.UpdateResponse
	UpgradeStateResponse   *resource.UpgradeStateResponse
	ValidateConfigResponse *resource.ValidateConfigResponse
	GenerateConfigResponse *resource.GenerateConfigResponse
}

func (r Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.CreateResponse != nil {
		resp.Diagnostics = r.CreateResponse.Diagnostics
		resp.NewState = r.CreateResponse.NewState
		resp.NewIdentity = r.CreateResponse.NewIdentity
	}

	if r.CreateFunc != nil {
		r.CreateFunc(ctx, req, resp)
	}
}

func (r Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.DeleteResponse != nil {
		resp.Diagnostics = r.DeleteResponse.Diagnostics
	}
}

func (r Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if r.ImportStateResponse != nil {
		resp.Diagnostics = r.ImportStateResponse.Diagnostics
		resp.State = r.ImportStateResponse.State
		resp.Identity = r.ImportStateResponse.Identity
	}
}

func (r Resource) PlanChange(ctx context.Context, req resource.PlanChangeRequest, resp *resource.PlanChangeResponse) {
	if r.PlanChangeFunc != nil {
		r.PlanChangeFunc(ctx, req, resp)
	}
}

func (r Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.ReadResponse != nil {
		resp.Diagnostics = r.ReadResponse.Diagnostics
		resp.NewState = r.ReadResponse.NewState
		resp.NewIdentity = r.ReadResponse.NewIdentity
	}
}

func (r Resource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	if r.IdentitySchemaResponse != nil {
		resp.Diagnostics = r.IdentitySchemaResponse.Diagnostics
		resp.Schema = r.IdentitySchemaResponse.Schema
	}
}

func (r Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	if r.SchemaResponse != nil {
		resp.Diagnostics = r.SchemaResponse.Diagnostics
		resp.Schema = r.SchemaResponse.Schema
	}
}

func (r Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.UpdateResponse != nil {
		resp.Diagnostics = r.UpdateResponse.Diagnostics
		resp.NewState = r.UpdateResponse.NewState
		resp.NewIdentity = r.UpdateResponse.NewIdentity
	}
}

func (r Resource) UpgradeState(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	if r.UpgradeStateResponse != nil {
		resp.Diagnostics = r.UpgradeStateResponse.Diagnostics
		resp.UpgradedState = r.UpgradeStateResponse.UpgradedState
	}
}

func (r Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if r.ValidateConfigResponse != nil {
		resp.Diagnostics = r.ValidateConfigResponse.Diagnostics
	}
}

func (r Resource) GenerateConfig(ctx context.Context, req resource.GenerateConfigRequest, resp *resource.GenerateConfigResponse) {
	if r.GenerateConfigResponse != nil {
		resp.GeneratedConfig = r.GenerateConfigResponse.GeneratedConfig
		resp.Diagnostics = r.GenerateConfigResponse.Diagnostics
	}
}
