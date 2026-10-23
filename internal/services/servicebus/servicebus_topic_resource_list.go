// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServiceBusTopicListResource struct{}

type ServiceBusTopicListModel struct {
	NamespaceId types.String `tfsdk:"servicebus_namespace_id"`
}

var _ sdk.FrameworkListWrappedResource = new(ServiceBusTopicListResource)

func (ServiceBusTopicListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = serviceBusTopicResourceName
}

func (ServiceBusTopicListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceServiceBusTopic()
}

func (ServiceBusTopicListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"servicebus_namespace_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{Func: namespaces.ValidateNamespaceID},
				},
			},
		},
	}
}

func (ServiceBusTopicListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ServiceBus.TopicsClient
	namespacesClient := metadata.Client.ServiceBus.NamespacesClient

	var data ServiceBusTopicListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentID, err := namespaces.ParseNamespaceID(data.NamespaceId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing Namespace ID for `%s`", serviceBusTopicResourceName), err)
		return
	}

	topicNamespaceID := topics.NewNamespaceID(parentID.SubscriptionId, parentID.ResourceGroupName, parentID.NamespaceName)

	resp, err := client.ListByNamespaceComplete(ctx, topicNamespaceID, topics.DefaultListByNamespaceOperationOptions())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", serviceBusTopicResourceName), err)
		return
	}

	// The List wrapper cancels the supplied context before the iterator below runs,
	// so capture its deadline to rebuild a live context for the per-item flatten.
	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", fmt.Errorf("context had no deadline"))
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		listCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range resp.Items {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceServiceBusTopic().Data(&terraform.InstanceState{})
			id, err := topics.ParseTopicIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", serviceBusTopicResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceServiceBusTopicFlatten(listCtx, rd, id, &item, namespacesClient, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", serviceBusTopicResourceName), err)
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
