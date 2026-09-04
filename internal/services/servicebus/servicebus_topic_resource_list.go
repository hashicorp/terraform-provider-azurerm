// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2024-01-01/topics"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServiceBusTopicListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(ServiceBusTopicListResource)

func (ServiceBusTopicListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceServiceBusTopic()
}

func (ServiceBusTopicListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = serviceBusTopicResourceName
}

func (ServiceBusTopicListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.ServiceBus.TopicsClient
	namespacesClient := metadata.Client.ServiceBus.NamespacesClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	var namespaceItems []namespaces.SBNamespace

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := namespacesClient.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", serviceBusTopicResourceName), err)
			return
		}
		namespaceItems = resp.Items
	default:
		resp, err := namespacesClient.ListComplete(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", serviceBusTopicResourceName), err)
			return
		}
		namespaceItems = resp.Items
	}

	results := make([]topics.SBTopic, 0)
	for _, ns := range namespaceItems {
		namespaceId, err := namespaces.ParseNamespaceID(pointer.From(ns.Id))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing ServiceBus Namespace ID for `%s`", serviceBusTopicResourceName), err)
			return
		}

		resp, err := client.ListByNamespaceComplete(ctx, topics.NewNamespaceID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName), topics.DefaultListByNamespaceOperationOptions())
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", serviceBusTopicResourceName), err)
			return
		}
		results = append(results, resp.Items...)
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

		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			id, err := topics.ParseTopicIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing ID for `%s`", serviceBusTopicResourceName), err)
				return
			}

			rd := resourceServiceBusTopic().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceServiceBusTopicFlatten(listCtx, rd, id, &item, namespacesClient); err != nil {
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
