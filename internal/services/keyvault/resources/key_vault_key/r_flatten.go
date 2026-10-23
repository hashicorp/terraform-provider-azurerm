// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"encoding/base64"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
)

func flattenKeyVaultKeyOptions(input *[]string) []any {
	results := make([]any, 0, len(*input))

	for _, option := range *input {
		results = append(results, option)
	}

	return results
}

func flattenKeyVaultKeyRotationPolicy(input keys.KeyRotationPolicy) []any {
	if input.LifetimeActions == nil && input.Attributes == nil {
		return []any{}
	}

	policy := make(map[string]any)
	if input.Attributes != nil && input.Attributes.ExpiryTime != nil && *input.Attributes.ExpiryTime != "" {
		policy["expire_after"] = *input.Attributes.ExpiryTime
	}

	if input.LifetimeActions != nil {
		for _, ltAction := range *input.LifetimeActions {
			action := ltAction.Action
			trigger := ltAction.Trigger

			if action != nil && trigger != nil && action.Type != nil && strings.EqualFold(string(*action.Type), string(keys.ActionTypeNotify)) && trigger.TimeBeforeExpiry != nil && *trigger.TimeBeforeExpiry != "" {
				// Somehow a default is set after creation for notify_before_expiry
				// Submitting this set value in the next run will not work though..
				if policy["expire_after"] != nil {
					policy["notify_before_expiry"] = *trigger.TimeBeforeExpiry
				}
			}

			if action != nil && trigger != nil && action.Type != nil && strings.EqualFold(string(*action.Type), string(keys.ActionTypeRotate)) {
				autoRotation := make(map[string]any)
				autoRotation["time_after_creation"] = pointer.From(trigger.TimeAfterCreate)
				autoRotation["time_before_expiry"] = pointer.From(trigger.TimeBeforeExpiry)
				policy["automatic"] = []map[string]any{autoRotation}
			}
		}
	}

	if len(policy) == 0 {
		return []any{}
	}

	return []any{policy}
}

func flattenKeyVaultKeyReleasePolicy(input *keys.KeyReleasePolicy) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	data := ""
	if input.Data != nil {
		decoded, err := base64.RawURLEncoding.DecodeString(*input.Data)
		if err != nil {
			return nil, err
		}
		data = string(decoded)
	}

	return []any{
		map[string]any{
			"json":      data,
			"immutable": pointer.From(input.Immutable),
		},
	}, nil
}
