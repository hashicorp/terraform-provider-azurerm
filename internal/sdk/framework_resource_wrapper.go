package sdk

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/commonschema"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type FrameworkResourceWrapper struct {
	ResourceMetadata

	FrameworkWrappedResource

	Model interface{}
}

var _ resource.ResourceWithModifyPlan = &FrameworkResourceWrapper{}

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
	var model = r.FrameworkWrappedResource.ModelObject()

	if ok := r.ResourceMetadata.DecodeCreate(ctx, request, response, model); !ok {
		return
	}

	r.FrameworkWrappedResource.Create(ctx, request, response, r.ResourceMetadata, model)
	if response.Diagnostics.HasError() {
		return
	}

	r.ResourceMetadata.EncodeCreate(ctx, response, model)
}

func (r *FrameworkResourceWrapper) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	ctx, cancel := context.WithTimeout(ctx, r.ResourceMetadata.TimeoutRead)
	defer cancel()

	var state = r.FrameworkWrappedResource.ModelObject()

	if ok := r.ResourceMetadata.DecodeRead(ctx, request, response, state); !ok {
		return
	}

	r.FrameworkWrappedResource.Read(ctx, request, response, r.ResourceMetadata, state)

	r.ResourceMetadata.EncodeRead(ctx, response, state)
}

func (r *FrameworkResourceWrapper) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	ctx, cancel := context.WithTimeout(ctx, *r.ResourceMetadata.TimeoutUpdate*time.Minute)
	defer cancel()

	var plan = r.FrameworkWrappedResource.ModelObject()
	var state = r.FrameworkWrappedResource.ModelObject()

	if ok := r.ResourceMetadata.DecodeUpdate(ctx, request, response, &plan, &state); !ok {
		return
	}

	r.FrameworkWrappedResource.Update(ctx, request, response, r.ResourceMetadata, plan, state)

	r.ResourceMetadata.EncodeUpdate(ctx, response, state)
}

func (r *FrameworkResourceWrapper) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	ctx, cancel := context.WithTimeout(ctx, r.ResourceMetadata.TimeoutDelete)
	defer cancel()

	var state = r.FrameworkWrappedResource.ModelObject()
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
