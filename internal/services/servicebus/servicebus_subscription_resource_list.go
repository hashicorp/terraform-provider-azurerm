package servicebus

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/subscriptions"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServiceBusSubscriptionListResource struct{}

type ServiceBusSubscriptionListModel struct {
	TopicId types.String `tfsdk:"topic_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ServiceBusSubscriptionListResource)

func (ServiceBusSubscriptionListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = serviceBusSubscriptionResourceName
}

func (ServiceBusSubscriptionListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceServiceBusSubscription()
}

func (ServiceBusSubscriptionListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"topic_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: subscriptions.ValidateTopicID},
				},
			},
		},
	}
}

func (ServiceBusSubscriptionListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ServiceBus.SubscriptionsClient

	var data ServiceBusSubscriptionListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := subscriptions.ParseTopicID(data.TopicId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Topic ID for `%s`", serviceBusSubscriptionResourceName), err)
		return
	}

	resp, err := client.ListByTopicComplete(ctx, *parentID, subscriptions.DefaultListByTopicOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", serviceBusSubscriptionResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceServiceBusSubscription().Data(&terraform.InstanceState{})

			id, err := subscriptions.ParseSubscriptions2IDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", serviceBusSubscriptionResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceServiceBusSubscriptionFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", serviceBusSubscriptionResourceName), err)
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
