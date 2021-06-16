package kusto

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	msiparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	msivalidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func schemaIdentity() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(kusto.IdentityTypeSystemAssigned),
						string(kusto.IdentityTypeUserAssigned),
						string(kusto.IdentityTypeSystemAssignedUserAssigned),
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},

				"identity_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: msivalidate.UserAssignedIdentityID,
					},
				},

				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandIdentity(input []interface{}) (*kusto.Identity, error) {
	if len(input) == 0 || input[0] == nil {
		return &kusto.Identity{
			Type: kusto.IdentityTypeNone,
		}, nil
	}
	raw := input[0].(map[string]interface{})

	kustoIdentity := kusto.Identity{
		Type: kusto.IdentityType(raw["type"].(string)),
	}

	identityIdsRaw := raw["identity_ids"].(*pluginsdk.Set).List()
	identityIds := make(map[string]*kusto.IdentityUserAssignedIdentitiesValue)
	for _, v := range identityIdsRaw {
		identityIds[v.(string)] = &kusto.IdentityUserAssignedIdentitiesValue{}
	}

	if len(identityIds) > 0 {
		if kustoIdentity.Type != kusto.IdentityTypeUserAssigned && kustoIdentity.Type != kusto.IdentityTypeSystemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}

		kustoIdentity.UserAssignedIdentities = identityIds
	}

	return &kustoIdentity, nil
}

func flattenIdentity(input *kusto.Identity) ([]interface{}, error) {
	if input == nil || input.Type == kusto.IdentityTypeNone {
		return []interface{}{}, nil
	}

	identityIds := make([]string, 0)
	if input.UserAssignedIdentities != nil {
		for key := range input.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityID(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
	}

	principalID := ""
	if input.PrincipalID != nil {
		principalID = *input.PrincipalID
	}

	tenantID := ""
	if input.TenantID != nil {
		tenantID = *input.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}, nil
}

func expandTrustedExternalTenants(input []interface{}) *[]kusto.TrustedExternalTenant {
	output := make([]kusto.TrustedExternalTenant, 0)

	for _, v := range input {
		output = append(output, kusto.TrustedExternalTenant{
			Value: utils.String(v.(string)),
		})
	}

	return &output
}

func flattenTrustedExternalTenants(input *[]kusto.TrustedExternalTenant) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		if v.Value == nil {
			continue
		}

		output = append(output, *v.Value)
	}

	return output
}
