package network

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceNetworkInterfaceApplicationSecurityGroupAssociationUpgradeV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"application_security_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceNetworkInterfaceApplicationSecurityGroupAssociationUpgradeV0ToV1(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
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
