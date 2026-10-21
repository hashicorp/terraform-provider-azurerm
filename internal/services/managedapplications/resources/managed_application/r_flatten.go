// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_application

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
)

func flattenManagedApplicationPlan(input *applications.Plan) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	results = append(results, map[string]any{
		"name":           input.Name,
		"product":        input.Product,
		"publisher":      input.Publisher,
		"version":        input.Version,
		"promotion_code": pointer.From(input.PromotionCode),
	})

	return results
}

func flattenManagedApplicationOutputs(input *any) (map[string]any, error) {
	results := make(map[string]any)
	if input == nil {
		return results, nil
	}

	attrs := *input
	if _, ok := attrs.(map[string]any); ok {
		for k, val := range attrs.(map[string]any) {
			mapVal, ok := val.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("unexpected managed application output type: %+v", mapVal)
			}
			if mapVal != nil {
				v, ok := mapVal["value"]
				if !ok {
					return nil, fmt.Errorf("missing key 'value' in output map %+v", mapVal)
				}

				value, err := extractParameterOrOutputValue(v)
				if err != nil {
					return nil, fmt.Errorf("extracting outputs: %+v", err)
				}
				results[k] = value
			}
		}
	}

	return results, nil
}

func flattenManagedApplicationParameterValuesValueToString(input *any, localParameters map[string]any) (string, error) {
	if input == nil {
		return "", nil
	}

	attrs := *input
	if _, ok := attrs.(map[string]any); ok {
		for k, v := range attrs.(map[string]any) {
			if v != nil {
				delete(attrs.(map[string]any)[k].(map[string]any), "type")

				// Secure values are not returned, thus settings it with local value
				value := attrs.(map[string]any)[k].(map[string]any)
				if _, ok := value["value"]; !ok {
					value["value"] = ""
					if localParam, localParamOK := localParameters[k]; localParamOK {
						if _, oldValueStructTypeOK := localParam.(map[string]any); oldValueStructTypeOK {
							if localParamValue, localParamValueOK := localParam.(map[string]any)["value"]; localParamValueOK {
								value["value"] = localParamValue
							}
						}
					}
				}
			}
		}

		return compactParameterOrOutputValue(input)
	}

	return "", nil
}

func flattenManagedApplicationIdentity(input *applications.Identity) ([]any, error) {
	var config *identity.SystemAndUserAssignedMap

	if input != nil {
		identityType := identity.TypeNone
		if input.Type != nil {
			identityType = identity.Type(*input.Type)
		}

		config = &identity.SystemAndUserAssignedMap{
			Type:        identityType,
			IdentityIds: nil,
		}

		if input.PrincipalId != nil {
			config.PrincipalId = *input.PrincipalId
		}

		if input.TenantId != nil {
			config.TenantId = *input.TenantId
		}

		identityIds := make(map[string]identity.UserAssignedIdentityDetails)
		if input.UserAssignedIdentities != nil {
			for k, v := range *input.UserAssignedIdentities {
				identityIds[k] = identity.UserAssignedIdentityDetails{
					PrincipalId: v.PrincipalId,
				}
			}
		}

		config.IdentityIds = identityIds
	}

	result, err := identity.FlattenSystemAndUserAssignedMap(config)
	if err != nil {
		return nil, err
	}

	return *result, nil
}
