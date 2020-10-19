package managementgroup

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagementGroupSubscriptionAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagementGroupSubscriptionAssociationCreate,
		Read:   resourceArmManagementGroupSubscriptionAssociationRead,
		Delete: resourceArmManagementGroupSubscriptionAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

			"subscription_id": {
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
	subscriptionGuid := d.Get("subscription_id").(string)
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

	resourceId := fmt.Sprintf("%s|%s", managementGroupId, subscriptionId)
	log.Printf("[INFO] resourceId is %q", resourceId)

	d.SetId(resourceId)

	return resourceArmManagementGroupSubscriptionAssociationRead(d, meta)
}

func resourceArmManagementGroupSubscriptionAssociationRead(d *schema.ResourceData, meta interface{}) error {
	groupsClient := meta.(*clients.Client).ManagementGroups.GroupsClient // No GET on SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementGroupSubscriptionAssociationID(d.Id())
	if err != nil {
		return err
	}

	recurse := false // Only want immediate children
	resp, err := groupsClient.Get(ctx, id.ManagementGroupName, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		if utils.ResponseWasForbidden(resp.Response) || utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Management Group %q doesn't exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Management Group %q: %+v", d.Id(), err)
	}

	props := resp.Properties
	if props == nil {
		return fmt.Errorf("properties was nil for Management Group %q: %+v", d.Id(), err)
	}

	children := props.Children
	if children == nil {
		log.Printf("[INFO] Management Group %q doesn't have any children - removing association with %q from state", d.Id(), id.SubscriptionID)
		d.SetId("")
		return nil
	}

	exists := false
	if subscriptionIds := *children; subscriptionIds != nil {
		for _, subscriptionId := range subscriptionIds {
			if strings.EqualFold(id.SubscriptionID, *subscriptionId.Name) {
				exists = true
			}
		}
	}

	if !exists {
		log.Printf("[INFO] Management Group %q doesn't include %q - %q - removing association from state", d.Id(), err, id.SubscriptionID)
		d.SetId("")
		return nil
	}

	d.Set("management_group_id", id.ManagementGroupID)
	d.Set("subscription_id", id.SubscriptionID)

	return nil
}

func resourceArmManagementGroupSubscriptionAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementGroupSubscriptionAssociationID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] id object: %+v", id)

	_, err = parse.ManagementGroupID(id.ManagementGroupID)
	if err != nil {
		return fmt.Errorf("Error parsing management group name '%s': %+v", id.ManagementGroupID, err)
	}

	log.Printf("[INFO] management group %q <-> subscription %q association deletion.", id.ManagementGroupName, id.SubscriptionID)

	locks.ByID(id.ManagementGroupID)
	defer locks.UnlockByID(id.ManagementGroupID)
	locks.ByID(id.SubscriptionScopeID)
	defer locks.UnlockByID(id.SubscriptionScopeID)

	result, err := client.Delete(ctx, id.ManagementGroupName, id.SubscriptionID, managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("Error deleting association for subscription %q to management group %q: %+v", id.SubscriptionID, id.ManagementGroupName, err)
	}

	log.Printf("[INFO] management group %q <-> subscription %q association deleted: %+v", id.ManagementGroupName, id.SubscriptionID, result)

	return nil
}
