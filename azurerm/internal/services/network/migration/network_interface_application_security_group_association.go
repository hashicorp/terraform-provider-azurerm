package migration

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func NetworkInterfaceApplicationSecurityGroupAssociationV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    networkInterfaceApplicationSecurityGroupAssociationSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: networkInterfaceApplicationSecurityGroupAssociationUpgradeV0ToV1,
		Version: 0,
	}
}

func networkInterfaceApplicationSecurityGroupAssociationSchemaForV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"application_security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func networkInterfaceApplicationSecurityGroupAssociationUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	// after shipping support for this Resource Azure's since changed the behaviour to require that all IP Configurations
	// are connected to the same Application Security Group
	applicationSecurityGroupId := rawState["application_security_group_id"].(string)
	networkInterfaceId := rawState["network_interface_id"].(string)

	oldID := rawState["id"].(string)
	newID := fmt.Sprintf("%s|%s", networkInterfaceId, applicationSecurityGroupId)
	log.Printf("[DEBUG] Updating ID from %q to %q", oldID, newID)

	rawState["id"] = newID
	return rawState, nil
}
