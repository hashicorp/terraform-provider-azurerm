# Framework Adoption - Experimental

**WARNING:** This functionality is experimental, and not for use in the provider at this time. It is intended for maintainer experimentation to facilitate migration efforts for moving from `terraform-plugin-sdk` to `terraform-plugin-framework`.  This package is subject to removal or significant breaking change. Any PR submitted referencing/using this package/functionality will not be accepted, and will be closed.

## Resources

Example: (Working re-implementation of `azurerm_resource_group` as Framework, named `azurerm_fw_resource_group` to avoid collision)
```go
package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	resourcegroupsvalidate "github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk/frameworkhelpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/custompollers"
)

var _ sdk.FrameworkResource = &FWResourceGroupResource{}

type FWResourceGroupResource struct {
	sdk.ResourceMetadata
}

func (r *FWResourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type FWResourceGroupResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Location  types.String `tfsdk:"location"`
	ManagedBy types.String `tfsdk:"managed_by"`
	Tags      types.Map    `tfsdk:"tags"`
}

func (r *FWResourceGroupResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Defaults(request, response)
}

func NewFWResourceGroupResource() resource.Resource {
	return &FWResourceGroupResource{}
}

func (r *FWResourceGroupResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_fw_resource_group"
}

func (r *FWResourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "",
				MarkdownDescription: "",
				Validators: []validator.String{
					frameworkhelpers.WrappedStringValidator{
						Func: resourcegroupsvalidate.ValidateName,
					},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"location": frameworkhelpers.LocationAttribute(),

			"tags": schema.MapAttribute{
				ElementType:         basetypes.StringType{},
				Optional:            true,
				Description:         "",
				MarkdownDescription: "",
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
				},
			},

			"managed_by": schema.StringAttribute{
				Optional:    true,
				Description: "",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},

			"id": frameworkhelpers.IDAttribute(),
		},
	}
}

func (r *FWResourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client := r.Client.Resource.ResourceGroupsClient
	createContext, cancel := context.WithTimeout(ctx, r.TimeoutCreate)
	defer cancel()

	var data FWResourceGroupResourceModel

	if ok := r.DecodeCreate(createContext, req, resp, &data); !ok {
		return
	}

	id := commonids.NewResourceGroupID(r.SubscriptionId, data.Name.ValueString())

	existing, err := client.Get(createContext, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("checking for presence of existing resource group: %+v", err), err.Error())
			return
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		r.ResourceRequiresImport("azurerm_fw_resource_group", id, resp)
		return
	}

	payload := resourcegroups.ResourceGroup{
		Location: location.Normalize(data.Location.ValueString()),
	}

	if !data.ManagedBy.IsNull() {
		payload.ManagedBy = data.ManagedBy.ValueStringPointer()
	}

	tags, diags := frameworkhelpers.ExpandTags(data.Tags)
	if diags.HasError() {
		sdk.AppendResponseErrorDiagnostic(resp, diags)
		return
	}

	payload.Tags = tags

	if _, err = client.CreateOrUpdate(createContext, id, payload); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("creating %s", id), err.Error())
		return
	}

	data.ID = types.StringValue(id.ID())

	r.EncodeCreate(createContext, resp, &data)
}

func (r *FWResourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client := r.Client.Resource.ResourceGroupsClient
	readContext, cancel := context.WithTimeout(ctx, r.TimeoutRead)
	defer cancel()

	var state FWResourceGroupResourceModel

	if ok := r.DecodeRead(readContext, req, resp, &state); !ok {
		return
	}

	id, err := commonids.ParseResourceGroupID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing resource group ID", err)
		return
	}

	existing, err := client.Get(readContext, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			r.MarkAsGone(id, &resp.State, &resp.Diagnostics)
			return
		}

		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id.String()), err)
		return
	}

	if model := existing.Model; model != nil {
		state.Name = types.StringValue(id.ResourceGroupName)
		state.Location = types.StringValue(location.Normalize(model.Location))
		state.ManagedBy = types.StringPointerValue(model.ManagedBy)
		t, diags := frameworkhelpers.FlattenTags(model.Tags)
		if diags.HasError() {
			sdk.AppendResponseErrorDiagnostic(resp, diags)
			return
		}

		state.Tags = t
	}

	r.EncodeRead(readContext, resp, &state)
}

func (r *FWResourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client := r.Client.Resource.ResourceGroupsClient
	updateContext, cancel := context.WithTimeout(ctx, *r.TimeoutUpdate)
	defer cancel()

	var plan, state FWResourceGroupResourceModel

	if ok := r.DecodeUpdate(updateContext, req, resp, &plan, &state); !ok {
		return
	}

	id, err := commonids.ParseResourceGroupID(state.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "Parsing ID", err)
		return
	}

	existing, err := client.Get(updateContext, *id)
	if err != nil && existing.Model != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving %s:", id.String()), err)
		return
	}

	model := existing.Model

	update := resourcegroups.ResourceGroupPatchable{
		ManagedBy: model.ManagedBy,
		Name:      model.Name,
		Tags:      model.Tags,
	}

	if !plan.ManagedBy.Equal(state.ManagedBy) {
		update.ManagedBy = plan.ManagedBy.ValueStringPointer()
	}

	if !plan.Tags.Equal(state.Tags) {
		tags, diags := frameworkhelpers.ExpandTags(plan.Tags)
		if diags.HasError() {
			sdk.AppendResponseErrorDiagnostic(resp, diags)
			return
		}
		update.Tags = tags
	}

	if _, err = client.Update(updateContext, *id, update); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("updating %s", id.String()), err)
		return
	}

	r.EncodeUpdate(updateContext, resp, &state)
}

func (r *FWResourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client := r.Client.Resource.ResourceGroupsClient
	deleteContext, cancel := context.WithTimeout(ctx, r.TimeoutDelete)
	defer cancel()

	var data FWResourceGroupResourceModel

	if ok := r.DecodeDelete(deleteContext, req, resp, &data); !ok {
		return
	}

	id, err := commonids.ParseResourceGroupID(data.ID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing ID", err)
		return
	}

	if r.Features.ResourceGroup.PreventDeletionIfContainsResources {
		pollerType := custompollers.NewResourceGroupEmptyPoller(r.Client.Resource.ResourcesClient, *id)
		poller := pollers.NewPoller(pollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err = poller.PollUntilDone(deleteContext); err != nil {
			sdk.SetResponseErrorDiagnostic(resp, "polling error", err)
		}
	}

	if err = client.DeleteThenPoll(deleteContext, *id, resourcegroups.DefaultDeleteOperationOptions()); err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("deleting %s:", id.String()), err.Error())
	}
}


```

## Data Sources

Example: // TODO

```go
package someservicepackage
// TODO

```

## Provider Functions (Core >= 1.8)

See `internal/provider/function` for live/shipped examples.

## Ephemeral Resources (Core >= 1.10)

Example:

```go
package someazureservice

import (
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type MyEphemeralResource struct {
	sdk.EphemeralResourceMetadata
}

var _ sdk.EphemeralResource = MyEphemeralResource{}

func (m MyEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "azurerm_my_ephemeral_resource"
}

func (m MyEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ...
		},
		Blocks: map[string]schema.Block{
			// ...
		},
	}
}

func (m MyEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	m.Defaults(req, resp)
}

func (m MyEphemeralResource) Open(ctx context.Context, request ephemeral.OpenRequest, openResponse *ephemeral.OpenResponse) {
	client := m.Client.SomeAzureService.FooClient

	// TODO - example code for ephemeral resource
}

```
