// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key_rotation_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneKeysClient

			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			// setting a blank policy to delete the existing policy
			if _, err := client.UpdateKeyRotationPolicy(ctx, id.BaseUri(), id.KeyName, keyvault.KeyRotationPolicy{
				LifetimeActions: pointer.To([]keyvault.LifetimeActions{}),
				Attributes:      &keyvault.KeyRotationPolicyAttributes{},
			}); err != nil {
				return fmt.Errorf("deleting HSM Key Rotation Policy for Key %q: %v", id, err)
			}

			return nil
		},
	}
}
