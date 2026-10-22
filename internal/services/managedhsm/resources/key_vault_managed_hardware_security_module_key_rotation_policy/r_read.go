// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_hardware_security_module_key_rotation_policy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
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

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			resp, err := client.GetKeyRotationPolicy(ctx, id.BaseUri(), id.KeyName)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			schema := flattenKeyRotationPolicy(resp)
			schema.ManagedHSMKeyID = id.ID()

			return metadata.Encode(&schema)
		},
	}
}

func flattenKeyRotationPolicy(p keyvault.KeyRotationPolicy) MHSMKeyRotationPolicyResourceSchema {
	var res MHSMKeyRotationPolicyResourceSchema
	if p.Attributes != nil {
		res.ExpireAfter = pointer.From(p.Attributes.ExpiryTime)
	}

	if p.LifetimeActions != nil {
		for _, ltAction := range *p.LifetimeActions {
			action := ltAction.Action
			trigger := ltAction.Trigger

			if action != nil && trigger != nil {
				if strings.EqualFold(string(action.Type), string(keyvault.ActionTypeRotate)) {
					res.TimeAfterCreation = pointer.From(trigger.TimeAfterCreate)
					res.TimeBeforeExpiry = pointer.From(trigger.TimeBeforeExpiry)
				}
			}
		}
	}

	return res
}

type MHSMKeyRotationPolicyResourceSchema struct {
	ManagedHSMKeyID   string `tfschema:"managed_hsm_key_id"`
	ExpireAfter       string `tfschema:"expire_after"`
	TimeAfterCreation string `tfschema:"time_after_creation"`
	TimeBeforeExpiry  string `tfschema:"time_before_expiry"`
}
