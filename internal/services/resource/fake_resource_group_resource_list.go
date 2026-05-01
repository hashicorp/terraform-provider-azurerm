package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type (
	FakeResourceGroupListResource struct{}
	FakeResourceGroupListModel    struct {
		SubscriptionId types.String `tfsdk:"subscription_id"`
		Filter         types.String `tfsdk:"filter"`
	}
)

var _ sdk.FrameworkListWrappedResourceWithConfig = new(FakeResourceGroupListResource)

func (r FakeResourceGroupListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = fakeResourceGroupResourceName
}

func (r FakeResourceGroupListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceFakeResourceGroup()
}

func (r FakeResourceGroupListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"subscription_id": listschema.StringAttribute{
				Optional:    true,
				Description: "The ID of the subscription to query. Defaults to the value specified in the Provider Configuration.",
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.IsUUID,
					},
				},
			},

			"filter": listschema.StringAttribute{
				Optional:    true,
				Description: "A filter expression to filter the results by.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (r FakeResourceGroupListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Resource.FakeResourceGroupsClient

	var data FakeResourceGroupListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]resourcegroups.FakeResourceGroup, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	options := resourcegroups.DefaultListOperationOptions()
	if !data.Filter.IsNull() {
		options = resourcegroups.ListOperationOptions{
			Filter: data.Filter.ValueStringPointer(),
		}
	}

	resp, err := client.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID), options)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", fakeResourceGroupResourceName), err)
		return
	}
	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, group := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(group.Name)

			id, err := commonids.ParseFakeResourceGroupID(pointer.From(group.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Resource Group ID", err)
				return
			}

			rd := resourceFakeResourceGroup().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceFakeResourceGroupFlatten(rd, id, &group); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` Resource Data", fakeResourceGroupResourceName), err)
				return
			}

			sdk.EncodeListResult(ctx, rd, &result)
			if result.Diagnostics.HasError() {
				push(result)
				return
			}

			if !push(result) {
				return
			}
		}
	}
}
