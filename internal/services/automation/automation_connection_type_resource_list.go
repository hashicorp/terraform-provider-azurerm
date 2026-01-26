package automation

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/connectiontype"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationConnectionTypeListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(AutomationConnectionTypeListResource)

func (r AutomationConnectionTypeListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(AutomationConnectionTypeResource{})
}

func (r AutomationConnectionTypeListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = automationConnectionTypeResourceName
}

type AutomationConnectionTypeListModel struct {
	SubscriptionId        types.String `tfsdk:"subscription_id"`
	ResourceGroupName     types.String `tfsdk:"resource_group_name"`
	AutomationAccountName types.String `tfsdk:"automation_account_name"`
}

func (r AutomationConnectionTypeListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"automation_account_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},

			"resource_group_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: resourcegroups.ValidateName,
					},
				},
			},

			"subscription_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateSubscriptionID,
					},
				},
			},
		},
	}
}

func (r AutomationConnectionTypeListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Automation.ConnectionType

	var data AutomationConnectionTypeListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	results := make([]connectiontype.ConnectionType, 0)

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	automationAccountId := connectiontype.NewAutomationAccountID(subscriptionID, data.ResourceGroupName.ValueString(), data.AutomationAccountName.ValueString())
	resp, err := client.ListByAutomationAccountComplete(ctx, automationAccountId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", automationConnectionTypeResourceName), err)
		return
	}
	results = resp.Items

	stream.Results = func(push func(list.ListResult) bool) {
		for _, profile := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(profile.Name)

			id, err := connectiontype.ParseConnectionTypeID(pointer.From(profile.Id))
			if err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, "parsing Network Profile ID", err)
				return
			}

			rd := sdk.WrappedResource(AutomationConnectionTypeResource{}).Data(&terraform.InstanceState{})
			resourceMetaData := sdk.NewResourceMetaData(metadata.Client, rd)

			rd.SetId(id.ID())

			if err := resourceAutomationConnectionTypeFlatten(resourceMetaData, id, &profile); err != nil {
				sdk.SetListIteratorErrorDiagnostic(result, push, fmt.Sprintf("encoding `%s` resource data", automationConnectionTypeResourceName), err)
				return
			}

			sdk.EncodeListResource(ctx, rd, result, push)
		}
	}
}
