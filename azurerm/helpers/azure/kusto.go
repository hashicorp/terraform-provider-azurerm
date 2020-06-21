package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-02-15/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
)

func SchemaKustoIdentity() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(kusto.IdentityTypeNone),
						string(kusto.IdentityTypeSystemAssigned),
					}, true),
					DiffSuppressFunc: suppress.CaseDifference,
				},
				"principal_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"identity_ids": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func ExpandKustoIdentity(input []interface{}) *kusto.Identity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	identity := input[0].(map[string]interface{})
	identityType := kusto.IdentityType(identity["type"].(string))

	kustoIdentity := kusto.Identity{
		Type: identityType,
	}

	return &kustoIdentity
}

func FlattenKustoIdentity(input *kusto.Identity) []interface{} {
	if input == nil || input.Type == kusto.IdentityTypeNone {
		return []interface{}{}
	}

	identityIds := make([]string, 0)
	if input.UserAssignedIdentities != nil {
		for k := range input.UserAssignedIdentities {
			identityIds = append(identityIds, k)
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
	}
}
