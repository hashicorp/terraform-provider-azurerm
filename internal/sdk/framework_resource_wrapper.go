package sdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type FrameworkResourceWrapper struct {
	ResourceMetadata

	FrameworkWrappedResource

	Model interface{}
}

var _ resource.ResourceWithModifyPlan = &FrameworkResourceWrapper{}

var _ resource.ResourceWithIdentity = &FrameworkResourceWrapper{}

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
	ctx, cancel := context.WithTimeout(ctx, r.ResourceMetadata.TimeoutCreate)
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
	ctx, cancel := context.WithTimeout(ctx, r.ResourceMetadata.TimeoutRead)
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
	ctx, cancel := context.WithTimeout(ctx, *r.ResourceMetadata.TimeoutUpdate)
	defer cancel()

	plan := r.FrameworkWrappedResource.ModelObject()
	state := r.FrameworkWrappedResource.ModelObject()

	r.ResourceMetadata.DecodeUpdate(ctx, request, response, &plan, &state)
	if response.Diagnostics.HasError() {
		return
	}

	r.FrameworkWrappedResource.Update(ctx, request, response, r.ResourceMetadata, plan, state)

	r.ResourceMetadata.EncodeUpdate(ctx, response, state)
}

func (r *FrameworkResourceWrapper) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	ctx, cancel := context.WithTimeout(ctx, r.ResourceMetadata.TimeoutDelete)
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
	if _, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithConfigure); ok {
		r.FrameworkWrappedResource.(FrameworkWrappedResourceWithConfigure).Configure(ctx, request, response, r.ResourceMetadata)
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
	if _, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithConfigValidators); ok {
		return r.FrameworkWrappedResource.(FrameworkWrappedResourceWithConfigValidators).ConfigValidators(ctx)
	}

	return nil
}

func (r *FrameworkResourceWrapper) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if _, ok := r.FrameworkWrappedResource.(FrameworkWrappedResourceWithPlanModifier); ok {
		r.FrameworkWrappedResource.(FrameworkWrappedResourceWithPlanModifier).ModifyPlan(ctx, request, response, r.ResourceMetadata)
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

		t := identityType(idType)

		segments := id.Segments()
		numSegments := len(segments)
		for idx, segment := range segments {
			if segmentTypeSupported(segment.Type) {
				name := segmentName(segment, t, numSegments, idx)

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

		t := identityType(idType)

		segments := id.Segments()
		numSegments := len(segments)
		for idx, segment := range segments {
			if segmentTypeSupported(segment.Type) {
				name := segmentName(segment, t, numSegments, idx)

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

func AssertResourceModelType[T any](input interface{}, response interface{}) *T {
	result, ok := input.(*T)
	if !ok {
		switch v := response.(type) {
		case *resource.CreateResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", input))
		case *resource.ReadResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", input))
		case *resource.UpdateResponse:
			v.Diagnostics.AddError("resource had wrong model type, ", fmt.Sprintf("got %T", input))
		case *resource.DeleteResponse:
			v.Diagnostics.AddError("resource had wrong model type", fmt.Sprintf("got %T", input))
		}

		return nil
	}

	return result
}
