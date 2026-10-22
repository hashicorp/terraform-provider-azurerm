// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneKeysClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultMHSMKeyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedHsmId, err := managedhsms.ParseManagedHSMID(config.ManagedHSMID)
			if err != nil {
				return err
			}
			baseUri, err := metadata.Client.ManagedHSMs.BaseUriForManagedHSM(ctx, *managedHsmId)
			if err != nil {
				return fmt.Errorf("determining the Data Plane Endpoint for %s: %+v", *managedHsmId, err)
			}
			if baseUri == nil {
				return fmt.Errorf("unable to determine the Data Plane Endpoint for %q", *managedHsmId)
			}
			endpoint, err := parse.ManagedHSMEndpoint(*baseUri, domainSuffix)
			if err != nil {
				return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", *endpoint, err)
			}

			id := parse.NewManagedHSMDataPlaneVersionlessKeyID(endpoint.ManagedHSMName, endpoint.DomainSuffix, config.Name)

			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetKey(ctx, endpoint.BaseURI(), id.KeyName, "")
				if err != nil {
					if !response.WasNotFound(existing.Response.Response) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.Response.Response) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := keyvault.KeyCreateParameters{
				Kty:    keyvault.JSONWebKeyType(config.KeyType),
				KeyOps: expandKeyVaultKeyOptions(config.KeyOpts),
				KeyAttributes: &keyvault.KeyAttributes{
					Enabled: pointer.To(true),
				},

				Tags: tags.Expand(config.Tags),
			}

			if config.Curve != "" {
				if config.KeyType != string(keyvault.JSONWebKeyTypeECHSM) {
					return fmt.Errorf("`key_type` must be `EC-HSM` when `curve` is set")
				}
				parameters.Curve = keyvault.JSONWebKeyCurveName(config.Curve)
			}

			if config.KeySize > 0 {
				if config.KeyType != string(keyvault.JSONWebKeyTypeRSAHSM) && config.KeyType != string(keyvault.JSONWebKeyTypeOctHSM) {
					return fmt.Errorf("`key_type` must be `RSA-HSM` or `oct-HSM` when `key_size` is set")
				}
				parameters.KeySize = pointer.To(int32(config.KeySize))
			}

			if config.NotBeforeDate != "" {
				notBeforeDate, _ := time.Parse(time.RFC3339, config.NotBeforeDate) // validated by schema
				parameters.KeyAttributes.NotBefore = pointer.To(date.UnixTime(notBeforeDate))
			}

			if config.ExpirationDate != "" {
				expirationDate, _ := time.Parse(time.RFC3339, config.ExpirationDate) // validated by schema
				parameters.KeyAttributes.Expires = pointer.To(date.UnixTime(expirationDate))
			}

			if resp, err := client.CreateKey(ctx, endpoint.BaseURI(), config.Name, parameters); err != nil {
				if metadata.Client.Features.KeyVault.RecoverSoftDeletedHSMKeys && response.WasConflict(resp.Response.Response) {
					recoveredKey, err := client.RecoverDeletedKey(ctx, endpoint.BaseURI(), config.Name)
					if err != nil {
						return err
					}
					log.Printf("[DEBUG] Recovering HSM Key %q with ID: %q", config.Name, *recoveredKey.Key.Kid)
					if kid := recoveredKey.Key.Kid; kid != nil {
						stateConf := &pluginsdk.StateChangeConf{
							Pending:                   []string{"pending"},
							Target:                    []string{"available"},
							Refresh:                   managedHSMKeyRefreshFunc(ctx, *kid),
							Delay:                     30 * time.Second,
							PollInterval:              10 * time.Second,
							ContinuousTargetOccurence: 10,
							Timeout:                   metadata.ResourceData.Timeout(pluginsdk.TimeoutCreate),
						}

						if _, err := stateConf.WaitForStateContext(ctx); err != nil {
							return fmt.Errorf("waiting for HSM Key %q to become available: %s", config.Name, err)
						}
						log.Printf("[DEBUG] Key %q recovered with ID: %q", config.Name, *kid)
					}
				} else {
					return fmt.Errorf("creating Key: %+v", err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandKeyVaultKeyOptions(input []string) *[]keyvault.JSONWebKeyOperation {
	results := make([]keyvault.JSONWebKeyOperation, 0, len(input))

	for _, option := range input {
		results = append(results, keyvault.JSONWebKeyOperation(option))
	}

	return &results
}
