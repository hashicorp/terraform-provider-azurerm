// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/list"

	"github.com/hashicorp/go-azure-helpers/framework/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type FrameworkResourceWrapper struct {
	ResourceMetadata

	FrameworkWrappedResource

	Model any
}

var _ resource.ResourceWithModifyPlan = &FrameworkResourceWrapper{}

var _ resource.ResourceWithIdentity = &FrameworkResourceWrapper{}

var _ list.ListResource = &FrameworkResourceWrapper{}

type EmbeddedFrameworkResourceModel interface{}

func (r *FrameworkResourceWrapper) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = r.FrameworkWrappedResource.ResourceType()
}

func (r *FrameworkResourceWrapper) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	r.FrameworkWrappedResource.Schema(ctx, request, response)
	response.Schema.Attributes["id"] = commonschema.IDAttribute()

	if response.Schema.Blocks == nil {
		response.Schema.Blocks = map[string]schema.Block{}
	}

	_, hasUpdate := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithUpdate)

	response.Schema.Blocks["timeouts"] = timeouts.Block(ctx, timeouts.Opts{
		Create: true,
		Read:   true,
		Update: hasUpdate,
		Delete: true,
	})
}

func (r *FrameworkResourceWrapper) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	customTimeouts := timeouts.Value{}
	response.Diagnostics.Append(request.Config.GetAttribute(ctx, path.Root("timeouts"), &customTimeouts)...)
	if response.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := customTimeouts.Create(ctx, r.ResourceMetadata.TimeoutCreate)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()
	model := r.FrameworkWrappedResource.ModelObject()

	r.ResourceMetadata.DecodeCreate(ctx, request, response, model)
	if response.Diagnostics.HasError() {
		return
	}

	r.FrameworkWrappedResource.Create(ctx, request, response, r.ResourceMetadata, model)
	if response.Diagnostics.HasError() {
		return
	}

	r.ResourceMetadata.EncodeCreate(ctx, response, model)
	if response.Diagnostics.HasError() {
		return
	}

	// Set the identity attributes on the response based on the resource ID encoded in the previous step.
	r.SetIdentityOnCreate(ctx, response)
}

func (r *FrameworkResourceWrapper) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	customTimeouts := timeouts.Value{}
	response.Diagnostics.Append(request.State.GetAttribute(ctx, path.Root("timeouts"), &customTimeouts)...)
	if response.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := customTimeouts.Read(ctx, r.ResourceMetadata.TimeoutRead)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	state := r.FrameworkWrappedResource.ModelObject()

	r.ResourceMetadata.DecodeRead(ctx, request, response, state)

	if response.Diagnostics.HasError() {
		return
	}

	r.FrameworkWrappedResource.Read(ctx, request, response, r.ResourceMetadata, state)
	if response.Diagnostics.HasError() {
		return
	}

	r.ResourceMetadata.EncodeRead(ctx, response, state)
	if response.Diagnostics.HasError() {
		return
	}

	r.SetIdentityOnRead(ctx, response)
}

func (r *FrameworkResourceWrapper) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	if fr, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithUpdate); ok {
		customTimeouts := timeouts.Value{}
		response.Diagnostics.Append(request.Config.GetAttribute(ctx, path.Root("timeouts"), &customTimeouts)...)
		if response.Diagnostics.HasError() {
			return
		}

		updateTimeout, diags := customTimeouts.Update(ctx, *r.ResourceMetadata.TimeoutUpdate)
		if diags.HasError() {
			response.Diagnostics.Append(diags...)
			return
		}

		ctx, cancel := context.WithTimeout(ctx, updateTimeout)
		defer cancel()

		plan := r.FrameworkWrappedResource.ModelObject()
		state := r.FrameworkWrappedResource.ModelObject()

		r.ResourceMetadata.DecodeUpdate(ctx, request, response, plan, state)
		if response.Diagnostics.HasError() {
			return
		}

		fr.Update(ctx, request, response, r.ResourceMetadata, plan, state)

		r.ResourceMetadata.EncodeUpdate(ctx, response, plan)

		return
	} else {
		SetResponseErrorDiagnostic(response, "Update called on non-updatable resource", fmt.Sprintf("resource type %s does not implement Update", r.FrameworkWrappedResource.ResourceType()))
	}
}

func (r *FrameworkResourceWrapper) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	customTimeouts := timeouts.Value{}
	response.Diagnostics.Append(request.State.GetAttribute(ctx, path.Root("timeouts"), &customTimeouts)...)
	if response.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := customTimeouts.Delete(ctx, r.ResourceMetadata.TimeoutDelete)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	state := r.FrameworkWrappedResource.ModelObject()
	r.ResourceMetadata.DecodeDelete(ctx, request, response, state)
	if response.Diagnostics.HasError() {
		return
	}

	r.FrameworkWrappedResource.Delete(ctx, request, response, r.ResourceMetadata, state)
}

func (r *FrameworkResourceWrapper) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Defaults(request, response)

	r.Model = r.ModelObject()
	if f, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithConfigure); ok {
		f.Configure(ctx, request, response, r.ResourceMetadata)
	}
}

func (r *FrameworkResourceWrapper) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	if request.Identity != nil {
		response.Identity = request.Identity
	}

	r.FrameworkWrappedResource.ImportState(ctx, request, response, r.ResourceMetadata)
}

func (r *FrameworkResourceWrapper) Resource() func() resource.Resource {
	return func() resource.Resource {
		return r
	}
}

func (r *FrameworkResourceWrapper) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if f, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithConfigValidators); ok {
		return f.ConfigValidators(ctx)
	}

	return nil
}

func (r *FrameworkResourceWrapper) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if f, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithPlanModifier); ok {
		f.ModifyPlan(ctx, request, response, r.ResourceMetadata)
	}
}

func (r *FrameworkResourceWrapper) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, response *resource.IdentitySchemaResponse) {
	response.IdentitySchema = GenerateIdentitySchema(r.FrameworkWrappedResource.Identity())
}

// SetIdentityOnCreate sets the identity attributes on the response based on the resource ID.
func (r *FrameworkResourceWrapper) SetIdentityOnCreate(ctx context.Context, response *resource.CreateResponse) {
	if id, idType := r.FrameworkWrappedResource.Identity(); id != nil {
		parser := resourceids.NewParserFromResourceIdType(id)
		idVal := ""
		response.State.GetAttribute(ctx, path.Root("id"), &idVal)
		parsed, err := parser.Parse(idVal, true)
		if err != nil {
			response.Diagnostics.AddError("parsing resource ID: %s", err.Error())
		}

		segments := id.Segments()
		numSegments := len(segments)
		for idx, segment := range segments {
			if segmentTypeSupported(segment.Type) {
				name := segmentName(segment, idType, numSegments, idx)

				field, ok := parsed.Parsed[segment.Name]
				if !ok {
					response.Diagnostics.AddError("setting resource identity", fmt.Sprintf("field `%s` was not found in the parsed resource ID %s", name, id))
					return
				}

				response.Identity.SetAttribute(ctx, path.Root(name), basetypes.NewStringValue(field))
			}
		}
	}
}

// SetIdentityOnRead sets the identity on the read response based on the resource ID.
func (r *FrameworkResourceWrapper) SetIdentityOnRead(ctx context.Context, response *resource.ReadResponse) {
	if id, idType := r.FrameworkWrappedResource.Identity(); id != nil {
		parser := resourceids.NewParserFromResourceIdType(id)
		idVal := ""
		response.State.GetAttribute(ctx, path.Root("id"), &idVal)
		parsed, err := parser.Parse(idVal, true)
		if err != nil {
			response.Diagnostics.AddError("parsing resource ID: %s", err.Error())
		}

		segments := id.Segments()
		numSegments := len(segments)
		for idx, segment := range segments {
			if segmentTypeSupported(segment.Type) {
				name := segmentName(segment, idType, numSegments, idx)

				field, ok := parsed.Parsed[segment.Name]
				if !ok {
					response.Diagnostics.AddError("setting resource identity", fmt.Sprintf("field `%s` was not found in the parsed resource ID %s", name, id))
					return
				}

				response.Identity.SetAttribute(ctx, path.Root(name), basetypes.NewStringValue(field))
			}
		}
	}
}

// List calls the supporting resource's List function to return the stream of results for the search query for that
// resource type.
func (r *FrameworkResourceWrapper) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream) {
	if lr, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithList); ok {
		lr.List(ctx, request, stream, r.ResourceMetadata)
	}
}

// ListResourceConfigSchema calls the supporting resource's ListResourceConfigSchema to populate the list response with
// the List Schema for the resource.
// If the resource does not implement List functionality, it will write an error diagnostic to inform the user that the
// resource does not (yet?) support List / Search.
func (r *FrameworkResourceWrapper) ListResourceConfigSchema(ctx context.Context, request list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	if lr, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithList); ok {
		lr.ListResourceConfigSchema(ctx, request, response, r.ResourceMetadata)
		return
	}

	response.Diagnostics.AddError("resource does not support list", fmt.Sprintf("the resource type %s does not support list/search", r.FrameworkWrappedResource.ResourceType()))
}

// AssertResourceModelType is a helper function to assist in checking the Resource or Data Source model type and
// provides the compiler with the struct type to be able to access the struct fields correctly.
// model is the initialised struct, and response is the operation response for which the model is being checked to
// allow error diagnostics to be written back to the SDK and passed to Terraform.
func AssertResourceModelType[T any](model any, response any) *T {
	result, ok := model.(*T)
	if !ok {
		switch v := response.(type) {
		case *resource.CreateResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", model))
		case *resource.ReadResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", model))
		case *resource.UpdateResponse:
			v.Diagnostics.AddError("resource had wrong model type, ", fmt.Sprintf("got %T", model))
		case *resource.DeleteResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", model))

		case *datasource.ReadResponse:
			v.Diagnostics.AddError("data source had wrong model type", fmt.Sprintf("got %T", model))
		}

		return nil
	}

	return result
}
