// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountConnectionApiKeyListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CognitiveAccountConnectionApiKeyListResource)

func (CognitiveAccountConnectionApiKeyListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CognitiveAccountConnectionApiKeyResource{})
}

func (CognitiveAccountConnectionApiKeyListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CognitiveAccountConnectionApiKeyResource{}.ResourceType()
}

func (CognitiveAccountConnectionApiKeyListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = cognitiveAccountConnectionListResourceConfigSchema()
}

func (CognitiveAccountConnectionApiKeyListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.AccountConnectionResourceClient

	var data cognitiveAccountConnectionListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	accountId, err := accountconnectionresource.ParseAccountID(data.CognitiveAccountId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Cognitive Account ID", err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		listCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		connectionsResp, err := client.AccountConnectionsListComplete(listCtx, accountconnectionresource.NewAccountID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName), accountconnectionresource.DefaultAccountConnectionsListOperationOptions())
		if err != nil {
			result := request.NewListResult(listCtx)
			sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("listing connections for `%s`", accountId.AccountName), err)
			return
		}

		for _, connection := range connectionsResp.Items {
			if connection.Properties == nil {
				continue
			}

			base := connection.Properties.ConnectionPropertiesV2()
			if base.AuthType != accountconnectionresource.ConnectionAuthTypeApiKey {
				continue
			}

			connectionId, err := accountconnectionresource.ParseConnectionID(pointer.From(connection.Id))
			if err != nil {
				result := request.NewListResult(listCtx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Account Connection ID", err)
				return
			}

			result := request.NewListResult(listCtx)
			result.DisplayName = pointer.From(connection.Name)

			r := CognitiveAccountConnectionApiKeyResource{}
			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(connectionId)

			if err := r.flatten(meta, connectionId, &connection, nil, ""); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", CognitiveAccountConnectionApiKeyResource{}.ResourceType()), err)
				return
			}

			sdk.EncodeListResult(listCtx, meta.ResourceData, &result)
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
