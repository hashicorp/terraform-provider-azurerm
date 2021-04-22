package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = NetworkInterfaceApplicationSecurityGroupAssociationV0ToV1{}

type NetworkInterfaceApplicationSecurityGroupAssociationV0ToV1 struct{}

func (NetworkInterfaceApplicationSecurityGroupAssociationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
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
	}
}

func (NetworkInterfaceApplicationSecurityGroupAssociationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
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
}
