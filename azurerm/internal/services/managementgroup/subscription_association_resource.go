package managementgroup

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	subscriptionValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceManagementGroupSubscriptionAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagementGroupSubscriptionAssociationCreate,
		Read:   resourceManagementGroupSubscriptionAssociationRead,
		Delete: resourceManagementGroupSubscriptionAssociationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: nil, // TODO

		Schema: map[string]*schema.Schema{
			"management_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagementGroupID,
			},

			"subscription_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: subscriptionValidate.SubscriptionID,
			},
		},
	}
}

func resourceManagementGroupSubscriptionAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId, err := parse.ManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}

	subscriptionId := d.Get("subscription_id").(string)

	_, err = client.Create(ctx, managementGroupId.Name, subscriptionId, "")
	if err != nil {
		return fmt.Errorf("creating Management Group Subscription Association between %q and %q: %+v", managementGroupId.Name, subscriptionId, err)
	}

	d.SetId(fmt.Sprintf("/managementGroup/%s/subscription/%s", managementGroupId.Name, subscriptionId))

	return nil
}

func resourceManagementGroupSubscriptionAssociationRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceManagementGroupSubscriptionAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
