package managementgroup

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagementGroupSubscriptionAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagementGroupSubscriptionAssociationCreate,
		Read:   resourceArmManagementGroupSubscriptionAssociationRead,
		Delete: resourceArmManagementGroupSubscriptionAssociationDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ManagementGroupID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"management_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagementGroupID,
			},

			"subscription_guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SubscriptionGUID,
			},
		},
	}
}

func resourceArmManagementGroupSubscriptionAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId := d.Get("management_group_id").(string)
	subscriptionGuid := d.Get("subscription_guid").(string)
	subscriptionId := "/subscriptions/" + subscriptionGuid

	parsedManagementGroupId, err := parse.ManagementGroupID(managementGroupId)
	if err != nil {
		return fmt.Errorf("Error parsing management group name '%s': %+v", managementGroupId, err)
	}

	managementGroupName := parsedManagementGroupId.Name

	log.Printf("[INFO] management group %q <-> subscription %q association creation.", managementGroupName, subscriptionGuid)

	locks.ByID(managementGroupId)
	defer locks.UnlockByID(managementGroupId)
	locks.ByID(subscriptionId)
	defer locks.UnlockByID(subscriptionId)

	result, err := client.Create(ctx, managementGroupName, subscriptionGuid, managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("Error creating association for subscription %q to management group %q: %+v", subscriptionGuid, managementGroupName, err)
	}

	log.Printf("[INFO] management group %q <-> subscription %q association created: %+v", managementGroupName, subscriptionGuid, result)

	id := fmt.Sprintf("%s|%s", managementGroupId, subscriptionGuid)

	d.SetId(id)

	d.Set("management_group_id", managementGroupId)
	d.Set("subscription_guid", subscriptionGuid)

	return resourceArmManagementGroupRead(d, meta)
}

func resourceArmManagementGroupSubscriptionAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient // No GET on SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId := d.Get("management_group_id").(string)
	subscriptionGuid := d.Get("subscription_guid").(string)

	parsedManagementGroupId, err := parse.ManagementGroupID(managementGroupId)
	if err != nil {
		return fmt.Errorf("Error parsing management group name '%s': %+v", managementGroupId, err)
	}

	managementGroupName := parsedManagementGroupId.Name

	recurse := false // Only want immediate children
	resp, err := client.Get(ctx, managementGroupName, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		if utils.ResponseWasForbidden(resp.Response) || utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Management Group %q doesn't exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Management Group %q: %+v", d.Id(), err)
	}

	if props := resp.Properties; props != nil {
		// subscriptionIds, err := flattenArmManagementGroupSubscriptionIds(props.Children)
		err := checkArmSubscriptionGuidInManagementGroupSubscriptionIds(subscriptionGuid, props.Children)
		if err != nil {
			log.Printf("[INFO] Management Group %q doesn't exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] Management Group Subscription association %q is not found - removing from state", d.Id())
	d.SetId("")
	return nil
}

func resourceArmManagementGroupSubscriptionAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId := d.Get("management_group_id").(string)
	subscriptionGuid := d.Get("subscription_guid").(string)
	subscriptionId := "/subscriptions/" + subscriptionGuid

	parsedManagementGroupId, err := parse.ManagementGroupID(managementGroupId)
	if err != nil {
		return fmt.Errorf("Error parsing management group name '%s': %+v", managementGroupId, err)
	}

	managementGroupName := parsedManagementGroupId.Name

	log.Printf("[INFO] management group %q <-> subscription %q association deletion.", managementGroupName, subscriptionGuid)

	locks.ByID(managementGroupId)
	defer locks.UnlockByID(managementGroupId)
	locks.ByID(subscriptionId)
	defer locks.UnlockByID(subscriptionId)

	result, err := client.Delete(ctx, managementGroupId, subscriptionGuid, managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("Error deleting association for subscription %q to management group %q: %+v", subscriptionGuid, managementGroupName, err)
	}

	log.Printf("[INFO] management group %q <-> subscription %q association deleted: %+v", managementGroupName, subscriptionGuid, result)

	return resourceArmManagementGroupRead(d, meta)
}

func checkArmSubscriptionGuidInManagementGroupSubscriptionIds(subscriptionGuid string, input *[]managementgroups.ChildInfo) error {
	if input == nil {
		return fmt.Errorf("no children")
	}

	for _, child := range *input {
		if child.ID == nil {
			continue
		}

		_, err := parseManagementGroupSubscriptionID(*child.ID)
		if err != nil {
			return fmt.Errorf("unable to parse child Subscription GUID %+v", err)
		}

		if *child.ID == subscriptionGuid {
			return nil
		}
	}

	return fmt.Errorf("%q was not found in children", subscriptionGuid)
}
