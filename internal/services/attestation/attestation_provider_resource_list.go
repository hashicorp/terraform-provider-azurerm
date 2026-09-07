// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package attestation

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AttestationProviderListResource struct{}

var _ sdk.FrameworkListWrappedResource = new(AttestationProviderListResource)

func (AttestationProviderListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = attestationProviderResourceName
}

func (AttestationProviderListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceAttestationProvider()
}

func (AttestationProviderListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.Attestation.ProviderClient

	var data sdk.DefaultListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	var results []attestationproviders.AttestationProviders
	subscriptionID := metadata.SubscriptionId
	if !data.SubscriptionId.IsNull() {
		subscriptionID = data.SubscriptionId.ValueString()
	}

	switch {
	case !data.ResourceGroupName.IsNull():
		resp, err := client.ListByResourceGroup(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", attestationProviderResourceName), err)
			return
		}
		if resp.Model != nil && resp.Model.Value != nil {
			results = *resp.Model.Value
		}
	default:
		resp, err := client.List(ctx, commonids.NewSubscriptionID(subscriptionID))
		if err != nil {
			sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("listing `%s`", attestationProviderResourceName), err)
			return
		}
		if resp.Model != nil && resp.Model.Value != nil {
			results = *resp.Model.Value
		}
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		sdk.SetResponseErrorDiagnostic(stream, "internal-error", "context had no deadline")
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		for _, item := range results {
			result := request.NewListResult(ctx)
			result.DisplayName = pointer.From(item.Name)

			rd := resourceAttestationProvider().Data(&terraform.InstanceState{})
			id, err := attestationproviders.ParseAttestationProvidersIDInsensitively(pointer.From(item.Id))
			if err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("parsing `%s` ID", attestationProviderResourceName), err)
				return
			}
			rd.SetId(id.ID())

			if err := resourceAttestationProviderFlatten(ctx, metadata.Client.Attestation, rd, id, &item, request.IncludeResource); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", attestationProviderResourceName), err)
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
