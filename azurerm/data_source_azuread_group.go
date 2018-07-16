package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmAzureADGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAzureADGroupRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"display_name"},
			},

			"display_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"object_id"},
			},

			"security_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceArmAzureADGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	var adgroup graphrbac.ADGroup
	var groupObj *graphrbac.ADGroup

	if oId, ok := d.GetOk("object_id"); ok {
		// use the object_id to find the Azure AD group
		resp, err := client.List(ctx, "")
		if err != nil {
			return fmt.Errorf("Error listing Azure AD groups: %+v", err)
		}

		for _, v := range *resp.Response().Value {
			if v.ObjectID != nil {
				if *v.ObjectID == oId {
					groupObj = &v
					break
				}
			}
		}
		if groupObj == nil {
			return fmt.Errorf("Couldn't locate a Azure AD group with the id %q", oId)
		}
	} else {
		// use the name to find the Azure AD group
		resp, err := client.ListComplete(ctx, "")
		if err != nil {
			return fmt.Errorf("Error listing Azure AD groups: %+v", err)
		}

		name := d.Get("display_name").(string)
		for _, v := range *resp.Response().Value {
			if v.DisplayName != nil {
				if *v.DisplayName == name {
					groupObj = &v
					break
				}
			}
		}
		if groupObj == nil {
			return fmt.Errorf("Couldn't locate a Azure AD group with a name of %q", name)
		}
	}

	adgroup = *groupObj

	d.SetId(*adgroup.ObjectID)
	d.Set("object_id", adgroup.ObjectID)
	d.Set("display_name", adgroup.DisplayName)
	d.Set("security_enabled", adgroup.SecurityEnabled)

	return nil
}
