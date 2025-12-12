// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// FrameworkListResourceWrapper presents an interface to simplify the implementations for List Resources
// it handles boilerplate code, and provides sensible defaults for the Azure Search APIs
type FrameworkListResourceWrapper struct {
	ResourceMetadata

	FrameworkListWrappedResource
}

type DefaultListModel struct {
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	SubscriptionId    types.String `tfsdk:"subscription_id"`
}

func (r *FrameworkListResourceWrapper) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	r.FrameworkListWrappedResource.Metadata(ctx, request, response)
}

func (r *FrameworkListResourceWrapper) ListResourceConfigSchema(ctx context.Context, request list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	if l, ok := r.FrameworkListWrappedResource.(FrameworkListWrappedResourceWithConfig); ok {
		l.ListResourceConfigSchema(ctx, request, response)
		return
	}

	// most resources default to RG and Subscription, so unless we need to customise that above, we can default it here.
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"resource_group_name": listschema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},
			"subscription_id": listschema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.IsUUID,
					},
				},
			},
		},
	}
}

func (r *FrameworkListResourceWrapper) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*60) // TODO - Custom Timeouts
	defer cancel()

	r.FrameworkListWrappedResource.List(ctx, request, stream, r.ResourceMetadata)
}

func (r *FrameworkListResourceWrapper) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, response *list.RawV5SchemaResponse) {
	res := r.FrameworkListWrappedResource.ResourceFunc()
	response.ProtoV5Schema = res.ProtoSchema(ctx)()
	response.ProtoV5IdentitySchema = res.ProtoIdentitySchema(ctx)()
}

func (r *FrameworkListResourceWrapper) Resource() func() list.ListResource {
	return func() list.ListResource {
		return r
	}
}

func (r *FrameworkListResourceWrapper) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Defaults(request, response)
}

var (
	_ list.ListResourceWithConfigure    = &FrameworkListResourceWrapper{}
	_ list.ListResourceWithRawV5Schemas = &FrameworkListResourceWrapper{}
)

type FrameworkListWrappedResource interface {
	Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse)

	List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata ResourceMetadata)

	// RawV5Schemas(ctx context.Context, request list.RawV5SchemaRequest, response *list.RawV5SchemaResponse)

	// ResourceFunc exposes the PluginSDKv2 resource function so that the RawV5Schema can be extracted
	// This should call the function that is used to register the resource in the provider that this List resource represents
	ResourceFunc() *pluginsdk.Resource
}

type FrameworkListWrappedResourceWithConfig interface {
	FrameworkListWrappedResource

	ListResourceConfigSchema(ctx context.Context, request list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse)
}
