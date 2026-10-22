// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key_rotation_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
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

			var config MHSMKeyRotationPolicyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			keyID, err := parse.ManagedHSMDataPlaneVersionlessKeyID(config.ManagedHSMKeyID, domainSuffix)
			if err != nil {
				return fmt.Errorf("parsing Managed HSM Key ID: %+v", err)
			}

			if _, err = client.GetKey(ctx, keyID.BaseUri(), keyID.KeyName, ""); err != nil {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", keyID, err)
			}

			// check key has rotation policy
			respPolicy, err := client.GetKeyRotationPolicy(ctx, keyID.BaseUri(), keyID.KeyName)
			if err != nil {
				switch {
				case response.WasForbidden(respPolicy.Response.Response):
					// If client is not authorized to access the policy:
					return fmt.Errorf("current client lacks permissions to read Key Rotation Policy for Key %q: %v", keyID, err)

				case response.WasNotFound(respPolicy.Response.Response):
					break
				default:
					return err
				}
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if respPolicy.Attributes != nil && respPolicy.Attributes.ExpiryTime != nil {
					if respPolicy.LifetimeActions != nil && len(*respPolicy.LifetimeActions) > 0 {
						return metadata.ResourceRequiresImport(r.ResourceType(), keyID)
					}
				}
			}

			if _, err := client.UpdateKeyRotationPolicy(ctx, keyID.BaseUri(), keyID.KeyName, expandKeyRotationPolicy(config)); err != nil {
				return fmt.Errorf("creating HSM Key Rotation Policy for Key %q: %v", keyID, err)
			}

			metadata.SetID(keyID)
			return nil
		},
	}
}

func expandKeyRotationPolicy(policy MHSMKeyRotationPolicyResourceSchema) keyvault.KeyRotationPolicy {
	var expiryTime *string // = nil // needs to be set to nil if not set
	if policy.ExpireAfter != "" {
		expiryTime = pointer.To(policy.ExpireAfter)
	}

	lifetimeActions := make([]keyvault.LifetimeActions, 0)

	lifetimeActionRotate := keyvault.LifetimeActions{
		Action: &keyvault.LifetimeActionsType{
			Type: keyvault.ActionTypeRotate,
		},
		Trigger: &keyvault.LifetimeActionsTrigger{},
	}

	if policy.TimeAfterCreation != "" {
		lifetimeActionRotate.Trigger.TimeAfterCreate = pointer.To(policy.TimeAfterCreation)
		lifetimeActions = append(lifetimeActions, lifetimeActionRotate)
	} else if policy.TimeBeforeExpiry != "" {
		lifetimeActionRotate.Trigger.TimeBeforeExpiry = pointer.To(policy.TimeBeforeExpiry)
		lifetimeActions = append(lifetimeActions, lifetimeActionRotate)
	}

	return keyvault.KeyRotationPolicy{
		LifetimeActions: &lifetimeActions,
		Attributes: &keyvault.KeyRotationPolicyAttributes{
			ExpiryTime: expiryTime,
		},
	}
}
