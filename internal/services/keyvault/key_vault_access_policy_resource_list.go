// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KeyVaultAccessPolicyListResource struct{}

type KeyVaultAccessPolicyListModel struct {
	KeyVaultId types.String `tfsdk:"key_vault_id"`
}

var _ sdk.FrameworkListWrappedResource = new(KeyVaultAccessPolicyListResource)

func (KeyVaultAccessPolicyListResource) ResourceFunc() *pluginsdk.Resource {
	return resourceKeyVaultAccessPolicy()
}

func (r KeyVaultAccessPolicyListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = keyVaultAccessPolicyResourceName
}

func (r KeyVaultAccessPolicyListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"key_vault_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateKeyVaultID,
					},
				},
			},
		},
	}
}

func (KeyVaultAccessPolicyListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {
	client := metadata.Client.KeyVault.VaultsClient

	var data KeyVaultAccessPolicyListModel
	diags := request.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	keyVaultId, err := commonids.ParseKeyVaultID(data.KeyVaultId.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("parsing `key_vault_id` for `%s`", keyVaultAccessPolicyResourceName), err)
		return
	}

	resp, err := client.Get(ctx, *keyVaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf("retrieving Key Vault for `%s`", keyVaultAccessPolicyResourceName), err)
		return
	}

	var accessPolicies []vaults.AccessPolicyEntry
	if resp.Model != nil && resp.Model.Properties.AccessPolicies != nil {
		accessPolicies = *resp.Model.Properties.AccessPolicies
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, policy := range accessPolicies {
			appId := ""
			if policy.ApplicationId != nil {
				appId = pointer.From(policy.ApplicationId)
			}

			id := parse.NewAccessPolicyId(*keyVaultId, policy.ObjectId, appId)

			result := request.NewListResult(ctx)
			result.DisplayName = fmt.Sprintf("%s/%s", keyVaultId.VaultName, policy.ObjectId)

			rd := resourceKeyVaultAccessPolicy().Data(&terraform.InstanceState{})
			rd.SetId(id.ID())

			if err := resourceKeyVaultAccessPolicyFlatten(rd, &id, &policy); err != nil {
				sdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf("encoding `%s` resource data", keyVaultAccessPolicyResourceName), err)
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
