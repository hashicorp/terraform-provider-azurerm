// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/rules"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServiceBusSubscriptionRuleListResource struct{}

type ServiceBusSubscriptionRuleListModel struct {
	ServicebusSubscriptionId types.String `tfsdk:"servicebus_subscription_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ServiceBusSubscriptionRuleListResource)

func (ServiceBusSubscriptionRuleListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = serviceBusSubscriptionRuleResourceName
}

func (ServiceBusSubscriptionRuleListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceServiceBusSubscriptionRule()
}

func (ServiceBusSubscriptionRuleListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"servicebus_subscription_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: rules.ValidateSubscriptions2ID},
				},
			},
		},
	}
}

func (ServiceBusSubscriptionRuleListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ServiceBus.SubscriptionRulesClient

	var data ServiceBusSubscriptionRuleListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := rules.ParseSubscriptions2ID(data.ServicebusSubscriptionId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Servicebus Subscription ID for `%s`", serviceBusSubscriptionRuleResourceName), err)
		return
	}

	resp, err := client.ListBySubscriptionsComplete(ctx, *parentID, rules.DefaultListBySubscriptionsOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", serviceBusSubscriptionRuleResourceName), err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)
			rd := resourceServiceBusSubscriptionRule().Data(&terraform.InstanceState{})
			id, err := rules.ParseRuleIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", serviceBusSubscriptionRuleResourceName), err)
				return
			}

			rd.SetId(id.ID())

			if err := resourceServiceBusSubscriptionRuleFlatten(rd, id, &item); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", serviceBusSubscriptionRuleResourceName), err)
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
