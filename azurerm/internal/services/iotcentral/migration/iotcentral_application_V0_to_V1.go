package migration

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/services/iotcentral/mgmt/2018-09-01/iotcentral"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func IoTCentralApplicationV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Version: 0,
		Type:    iotCentralApplicationV0Schema().CoreConfigSchema().ImpliedType(),
		Upgrade: iotCentralApplicationUpgradeV0ToV1,
	}
}

func iotCentralApplicationV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApplicationName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sub_domain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ApplicationSubdomain,
			},

			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ApplicationDisplayName,
			},

			"sku": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(iotcentral.F1),
					string(iotcentral.S1),
					string(iotcentral.ST1),
					string(iotcentral.ST2),
				}, true),
				Default: iotcentral.ST1,
			},
			"template": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ApplicationTemplateName,
			},

			"tags": tags.Schema(),
		},
	}
}

func iotCentralApplicationUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	id, err := parse.ApplicationID(oldId)
	if err != nil {
		return rawState, err
	}

	newId := id.ID()
	log.Printf("Updating `id` from %q to %q", oldId, newId)
	rawState["id"] = newId
	return rawState, nil
}
