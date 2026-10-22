// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneKeysClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			schema := KeyVaultMHSMKeyResourceSchema{}

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			resourceManagerId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseUri(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining Resource Manager ID for %q: %+v", id, err)
			}
			if resourceManagerId == nil {
				return metadata.MarkAsGone(*id)
			}

			resp, err := client.GetKey(ctx, id.BaseUri(), id.KeyName, "")
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if key := resp.Key; key != nil {
				schema.Name = id.KeyName
				schema.ManagedHSMID = resourceManagerId.ID()
				schema.KeyType = string(key.Kty)
				schema.KeyOpts = flattenKeyVaultKeyOptions(key.KeyOps)
				schema.Curve = string(key.Crv)
				schema.Tags = tags.Flatten(resp.Tags)
				schema.VersionedId = pointer.From(key.Kid)
				if key.N != nil {
					nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
					if err != nil {
						return fmt.Errorf("could not decode N: %+v", err)
					}
					schema.KeySize = int64(len(nBytes) * 8)
				}

				if attributes := resp.Attributes; attributes != nil {
					if v := attributes.NotBefore; v != nil {
						schema.NotBeforeDate = time.Time(*v).Format(time.RFC3339)
					}

					if v := attributes.Expires; v != nil {
						schema.ExpirationDate = time.Time(*v).Format(time.RFC3339)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func flattenKeyVaultKeyOptions(input *[]string) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	return append(results, *input...)
}

type KeyVaultMHSMKeyResourceSchema struct {
	Name           string         `tfschema:"name"`
	ManagedHSMID   string         `tfschema:"managed_hsm_id"`
	KeyType        string         `tfschema:"key_type"`
	KeyOpts        []string       `tfschema:"key_opts"`
	KeySize        int64          `tfschema:"key_size"`
	Curve          string         `tfschema:"curve"`
	NotBeforeDate  string         `tfschema:"not_before_date"`
	ExpirationDate string         `tfschema:"expiration_date"`
	Tags           map[string]any `tfschema:"tags"`
	VersionedId    string         `tfschema:"versioned_id"`
}

type deleteAndPurgeKey struct {
	client      *keyvault.BaseClient
	keyVaultUri string
	name        string
}

func (d deleteAndPurgeKey) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.DeleteKey(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func (d deleteAndPurgeKey) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetKey(ctx, d.keyVaultUri, d.name, "")
	return resp.Response, err
}

func (d deleteAndPurgeKey) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	return d.client.PurgeDeletedKey(ctx, d.keyVaultUri, d.name)
}

func (d deleteAndPurgeKey) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetDeletedKey(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}
