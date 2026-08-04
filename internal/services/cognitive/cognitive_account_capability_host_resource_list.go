// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountcapabilityhost"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountCapabilityHostListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(CognitiveAccountCapabilityHostListResource)

func (CognitiveAccountCapabilityHostListResource) ResourceFunc() *pluginsdk.Resource {
	return sdk.WrappedResource(CognitiveAccountCapabilityHostResource{})
}

func (CognitiveAccountCapabilityHostListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = CognitiveAccountCapabilityHostResource{}.ResourceType()
}

type cognitiveAccountCapabilityHostListModel struct {
	CognitiveAccountId types.String `tfsdk:"cognitive_account_id"`
}

func (CognitiveAccountCapabilityHostListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cognitive_account_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: cognitiveservicesaccounts.ValidateAccountID,
					},
				},
			},
		},
	}
}

func (CognitiveAccountCapabilityHostListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Cognitive.AccountCapabilityHostClient

	var data cognitiveAccountCapabilityHostListModel
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

	accountId, err := cognitiveservicesaccounts.ParseAccountID(data.CognitiveAccountId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, "parsing Cognitive Account ID", err)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		listCtx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		capabilityHostAccountId := accountcapabilityhost.NewAccountID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName)

		capabilityHostsResp, err := client.ListComplete(listCtx, capabilityHostAccountId)
		if err != nil {
			result := request.NewListResult(listCtx)
			sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("listing capability hosts for `%s`", accountId.AccountName), err)
			return
		}

		for _, capabilityHost := range capabilityHostsResp.Items {
			capabilityHostId, err := accountcapabilityhost.ParseCapabilityHostID(pointer.From(capabilityHost.Id))
			if err != nil {
				result := request.NewListResult(listCtx)
				sdk.SetErrorDiagnosticAndPushListResult(result, push, "parsing Cognitive Account Capability Host ID", err)
				return
			}

			result := request.NewListResult(listCtx)
			result.DisplayName = pointer.From(capabilityHost.Name)

			r := CognitiveAccountCapabilityHostResource{}
			meta := sdk.NewResourceMetaData(metadata.Client, r)
			meta.SetID(capabilityHostId)

			if err := r.flatten(meta, capabilityHostId, &capabilityHost); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", CognitiveAccountCapabilityHostResource{}.ResourceType()), err)
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
