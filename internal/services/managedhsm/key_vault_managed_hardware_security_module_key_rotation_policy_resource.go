// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	validate2 "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

type KeyVaultMHSMKeyRotationPolicyResource struct{}

var _ sdk.ResourceWithUpdate = KeyVaultMHSMKeyRotationPolicyResource{}

func (r KeyVaultMHSMKeyRotationPolicyResource) ModelObject() interface{} {
	return &MHSMKeyRotationPolicyResourceSchema{}
}

type RotateAutomatic struct {
	TimeAfterCreation string `tfschema:"time_after_creation"`
	TimeBeforeExpiry  string `tfschema:"time_before_expiry"`
}

type MHSMKeyRotationPolicyResourceSchema struct {
	ManagedHSMKeyID    string            `tfschema:"managed_hsm_key_id"`
	ExipreAfter        string            `tfschema:"expire_after"`
	NotifyBeforeExpiry string            `tfschema:"notify_before_expiry"`
	Automatic          []RotateAutomatic `tfschema:"automatic"`
}

func (r KeyVaultMHSMKeyRotationPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedHSMDataPlaneVersionlessKeyID
}

func (r KeyVaultMHSMKeyRotationPolicyResource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_key_rotation_policy"
}

func (r KeyVaultMHSMKeyRotationPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_hsm_key_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"expire_after": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate2.ISO8601DurationBetween("P28D", "P100Y"),
			AtLeastOneOf: []string{
				"expire_after",
				"automatic",
			},
			RequiredWith: []string{
				"expire_after",
				"notify_before_expiry",
			},
		},

		// <= expiry_time - 7, >=7
		"notify_before_expiry": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate2.ISO8601DurationBetween("P7D", "P36493D"),
			RequiredWith: []string{
				"expire_after",
				"notify_before_expiry",
			},
		},

		"automatic": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"time_after_creation": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate2.ISO8601Duration,
						AtLeastOneOf: []string{
							"automatic.0.time_after_creation",
							"automatic.0.time_before_expiry",
						},
					},
					"time_before_expiry": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate2.ISO8601Duration,
						AtLeastOneOf: []string{
							"automatic.0.time_after_creation",
							"automatic.0.time_before_expiry",
						},
					},
				},
			},
		},
	}
}

func (r KeyVaultMHSMKeyRotationPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"versioned_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r KeyVaultMHSMKeyRotationPolicyResource) Create() sdk.ResourceFunc {
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

			_, err = client.GetKey(ctx, keyID.BaseUri(), keyID.KeyName, "")
			if err != nil {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", keyID, err)
			}

			// check key has rotation policy
			respPolicy, err := client.GetKeyRotationPolicy(ctx, keyID.BaseUri(), keyID.KeyName)
			if err != nil {
				switch {
				case utils.ResponseWasForbidden(respPolicy.Response):
					// If client is not authorized to access the policy:
					return fmt.Errorf("current client lacks permissions to read Key Rotation Policy for Key %q: %v", keyID, err)

				case utils.ResponseWasNotFound(respPolicy.Response):
					break
				default:
					return err
				}
			}

			if respPolicy.Attributes != nil && respPolicy.Attributes.ExpiryTime != nil {
				if respPolicy.LifetimeActions != nil && len(*respPolicy.LifetimeActions) > 0 {
					return metadata.ResourceRequiresImport(r.ResourceType(), keyID)
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

func (r KeyVaultMHSMKeyRotationPolicyResource) Read() sdk.ResourceFunc {
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
				if utils.ResponseWasNotFound(resp.Response) {
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

func (r KeyVaultMHSMKeyRotationPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config MHSMKeyRotationPolicyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			if _, err := client.UpdateKeyRotationPolicy(ctx, id.BaseUri(), id.KeyName, expandKeyRotationPolicy(config)); err != nil {
				return fmt.Errorf("updating HSM Key Rotation Policy for Key %q: %v", id, err)
			}

			return nil
		},
	}
}

func (r KeyVaultMHSMKeyRotationPolicyResource) Delete() sdk.ResourceFunc {
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

			// settign a blank policy to delete the existing policy
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

func expandKeyRotationPolicy(policy MHSMKeyRotationPolicyResourceSchema) keyvault.KeyRotationPolicy {

	var expiryTime *string // = nil // needs to be set to nil if not set
	if policy.ExipreAfter != "" {
		expiryTime = utils.String(policy.ExipreAfter)
	}

	lifetimeActions := make([]keyvault.LifetimeActions, 0)
	if policy.NotifyBeforeExpiry != "" {
		lifetimeActionNotify := keyvault.LifetimeActions{
			Trigger: &keyvault.LifetimeActionsTrigger{
				TimeBeforeExpiry: utils.String(policy.NotifyBeforeExpiry), // for Type: keyvault.Notify always TimeBeforeExpiry
			},
			Action: &keyvault.LifetimeActionsType{
				Type: keyvault.ActionTypeNotify,
			},
		}
		lifetimeActions = append(lifetimeActions, lifetimeActionNotify)
	}

	if len(policy.Automatic) == 1 {
		lifetimeActionRotate := keyvault.LifetimeActions{
			Action: &keyvault.LifetimeActionsType{
				Type: keyvault.ActionTypeRotate,
			},
			Trigger: &keyvault.LifetimeActionsTrigger{},
		}
		autoItem := policy.Automatic[0]

		if autoItem.TimeAfterCreation != "" {
			lifetimeActionRotate.Trigger.TimeAfterCreate = pointer.To(autoItem.TimeAfterCreation)
		}

		if autoItem.TimeBeforeExpiry != "" {
			lifetimeActionRotate.Trigger.TimeBeforeExpiry = pointer.To(autoItem.TimeBeforeExpiry)
		}

		lifetimeActions = append(lifetimeActions, lifetimeActionRotate)
	}

	return keyvault.KeyRotationPolicy{
		LifetimeActions: &lifetimeActions,
		Attributes: &keyvault.KeyRotationPolicyAttributes{
			ExpiryTime: expiryTime,
		},
	}
}

func flattenKeyRotationPolicy(p keyvault.KeyRotationPolicy) MHSMKeyRotationPolicyResourceSchema {
	var res MHSMKeyRotationPolicyResourceSchema
	if p.Attributes != nil {
		res.ExipreAfter = pointer.From(p.Attributes.ExpiryTime)
	}

	if p.LifetimeActions != nil {
		for _, ltAction := range *p.LifetimeActions {
			action := ltAction.Action
			trigger := ltAction.Trigger

			if action != nil && trigger != nil {
				switch action.Type {
				case keyvault.ActionTypeNotify:
					res.NotifyBeforeExpiry = pointer.From(trigger.TimeBeforeExpiry)
				case keyvault.ActionTypeRotate:
					rotate := RotateAutomatic{
						TimeAfterCreation: pointer.From(trigger.TimeAfterCreate),
						TimeBeforeExpiry:  pointer.From(trigger.TimeBeforeExpiry),
					}
					res.Automatic = append(res.Automatic, rotate)
				}
			}
		}
	}

	return res
}
