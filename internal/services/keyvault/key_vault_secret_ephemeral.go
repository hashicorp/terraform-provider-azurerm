package keyvault

import (
	"context"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
	Name       types.String `tfsdk:"name"`
	KeyVaultID types.String `tfsdk:"key_vault_id"`
	Value      types.String `tfsdk:"value"`
	Version    types.String `tfsdk:"version"`
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
				Required:   true,
				Validators: nil, // TODO
			},

			"key_vault_id": schema.StringAttribute{
				Required:   true,
				Validators: nil, // TODO
			},

			"value": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},

			"version": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (e *KeyVaultSecretEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	keyVaultsClient := e.Client.KeyVault
	client := e.Client.KeyVault.ManagementClient

	var data KeyVaultSecretEphemeralResourceModel

	if ok := e.DecodeOpen(ctx, req, resp, &data); !ok {
		return
	}

	keyVaultId, err := commonids.ParseKeyVaultID(data.KeyVaultID.ValueString())
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "parsing", err)
		return
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		sdk.SetResponseErrorDiagnostic(resp, "looking up secret %q vault url from id", err)
	}

	response, err := client.GetSecret(ctx, *keyVaultBaseUri, data.Name.ValueString(), data.Version.ValueString())
	if err != nil {
		if utils.ResponseWasNotFound(response.Response) {
			sdk.SetResponseErrorDiagnostic(resp, "keyvault secret does not exist", err)
		}
		sdk.SetResponseErrorDiagnostic(resp, "retrieving", err)
	}

	data.Value = types.StringValue(pointer.From(response.Value))

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
