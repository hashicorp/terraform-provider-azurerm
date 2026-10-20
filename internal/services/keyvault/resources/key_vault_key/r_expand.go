// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"encoding/base64"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandKeyVaultKeyOptions(d *pluginsdk.ResourceData) *[]keys.JsonWebKeyOperation {
	options := d.Get("key_opts").([]any)
	results := make([]keys.JsonWebKeyOperation, 0, len(options))

	for _, option := range options {
		results = append(results, keys.JsonWebKeyOperation(option.(string)))
	}

	return &results
}

func expandKeyVaultKeyRotationPolicy(v []any) keys.KeyRotationPolicy {
	if len(v) == 0 {
		return keys.KeyRotationPolicy{LifetimeActions: &[]keys.LifetimeActions{}}
	}

	policy := v[0].(map[string]any)

	var expiryTime *string = nil // needs to be set to nil if not set
	if rawExpiryTime := policy["expire_after"]; rawExpiryTime != nil && rawExpiryTime.(string) != "" {
		expiryTime = pointer.To(rawExpiryTime.(string))
	}

	lifetimeActions := make([]keys.LifetimeActions, 0)
	if rawNotificationTime := policy["notify_before_expiry"]; rawNotificationTime != nil && rawNotificationTime.(string) != "" {
		lifetimeActionNotify := keys.LifetimeActions{
			Trigger: &keys.LifetimeActionsTrigger{
				TimeBeforeExpiry: pointer.To(rawNotificationTime.(string)), // for Type: Notify always TimeBeforeExpiry
			},
			Action: &keys.LifetimeActionsType{
				Type: pointer.To(keys.ActionTypeNotify),
			},
		}
		lifetimeActions = append(lifetimeActions, lifetimeActionNotify)
	}

	if autoRotationList := policy["automatic"].([]any); len(autoRotationList) == 1 && autoRotationList[0] != nil {
		lifetimeActionRotate := keys.LifetimeActions{
			Action: &keys.LifetimeActionsType{
				Type: pointer.To(keys.ActionTypeRotate),
			},
			Trigger: &keys.LifetimeActionsTrigger{},
		}
		autoRotationRaw := autoRotationList[0].(map[string]any)

		if v := autoRotationRaw["time_after_creation"]; v != nil && v.(string) != "" {
			lifetimeActionRotate.Trigger.TimeAfterCreate = pointer.To(v.(string))
		}

		if v := autoRotationRaw["time_before_expiry"]; v != nil && v.(string) != "" {
			lifetimeActionRotate.Trigger.TimeBeforeExpiry = pointer.To(v.(string))
		}

		lifetimeActions = append(lifetimeActions, lifetimeActionRotate)
	}

	return keys.KeyRotationPolicy{
		LifetimeActions: &lifetimeActions,
		Attributes: &keys.KeyRotationPolicyAttributes{
			ExpiryTime: expiryTime,
		},
	}
}

func expandKeyVaultKeyReleasePolicy(input []any) *keys.KeyReleasePolicy {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &keys.KeyReleasePolicy{
		Data:      pointer.To(base64.StdEncoding.EncodeToString([]byte(v["json"].(string)))),
		Immutable: pointer.To(v["immutable"].(bool)),
	}
}
