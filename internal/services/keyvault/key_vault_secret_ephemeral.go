// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/framework/typehelpers"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.EphemeralResource = &KeyVaultSecretEphemeralResource{}

func NewKeyVaultSecretEphemeralResource() ephemeral.EphemeralResource {
	return &KeyVaultSecretEphemeralResource{}
}

type KeyVaultSecretEphemeralResource struct {
	sdk.EphemeralResourceMetadata
}

type KeyVaultSecretEphemeralResourceModel struct {
	Name           types.String `tfsdk:"name"`
	KeyVaultID     types.String `tfsdk:"key_vault_id"`
	Version        types.String `tfsdk:"version"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	NotBeforeDate  types.String `tfsdk:"not_before_date"`
	Value          types.String `tfsdk:"value"`
}

func (e *KeyVaultSecretEphemeralResource) Metadata(_ context.Context, _ ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = "azurerm_key_vault_secret"
}

func (e *KeyVaultSecretEphemeralResource) Configure(_ context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	e.Defaults(req, resp)
}

func (e *KeyVaultSecretEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringIsNotEmpty,
					},
				},
			},

			"key_vault_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: commonids.ValidateKeyVaultID,
					},
				},
			},

			"version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					typehelpers.WrappedStringValidator{
						Func: validation.StringIsNotEmpty,
					},
				},
			},

			"expiration_date": schema.StringAttribute{
				Computed: true,
			},

			"not_before_date": schema.StringAttribute{
				Computed: true,
			},

			"value": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (e *KeyVaultSecretEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	keyVaultsClient := e.Client.KeyVault
	client := e.Client.KeyVault.ManagementClient
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	var data KeyVaultSecretEphemeralResourceModel

	if ok := e.DecodeOpen(ctx, req, resp, &data); !ok {
		return
	}

	keyVaultID, err := commonids.ParseKeyVaultID(data.KeyVaultID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "", err)
		return
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("looking up base uri for secret %q in %s", data.Name.ValueString(), keyVaultID), err)
		return
	}

	response, err := client.GetSecret(ctx, *keyVaultBaseUri, data.Name.ValueString(), data.Version.ValueString())
	if err != nil {
		if utils.ResponseWasNotFound(response.Response) {
			sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("secret %s does not exist in %s", data.Name.ValueString(), keyVaultID), err)
			return
		}
		sdk.SetResponseErrorDiagnostic(resp, fmt.Sprintf("retrieving secret %q from %s", data.Name.ValueString(), keyVaultID), err)
		return
	}

	data.Value = types.StringValue(pointer.From(response.Value))

	id, err := parse.ParseNestedItemID(*response.ID)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "", err)
		return
	}

	data.Version = types.StringValue(id.Version)

	if attributes := response.Attributes; attributes != nil {
		if expirationDate := attributes.Expires; expirationDate != nil {
			data.ExpirationDate = types.StringValue(time.Time(*expirationDate).Format(time.RFC3339))
		}

		if notBeforeDate := attributes.NotBefore; notBeforeDate != nil {
			data.NotBeforeDate = types.StringValue(time.Time(*notBeforeDate).Format(time.RFC3339))
		}
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
