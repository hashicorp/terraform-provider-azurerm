package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/management"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceManagementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagementGroupCreateUpdate,
		Update: resourceManagementGroupCreateUpdate,
		Read:   resourceManagementGroupRead,
		Delete: resourceManagementGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subscription_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceManagementGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient).managementGroupsClient
	subscriptionsClient := meta.(*ArmClient).managementGroupsSubscriptionClient
	ctx := meta.(*ArmClient).StopContext

	armTenantID := meta.(*ArmClient).tenantId
	name := d.Get("name").(string)
	subscriptionIds := d.Get("subscription_ids").([]interface{})
	log.Printf("[INFO] Creating management group %q", name)

	parentID := fmt.Sprintf("/providers/Microsoft.Management/managementGroups/%s", armTenantID)
	properties := managementgroups.CreateManagementGroupRequest{
		CreateManagementGroupProperties: &managementgroups.CreateManagementGroupProperties{
			TenantID:    &armTenantID,
			DisplayName: &name,
			Details: &managementgroups.CreateManagementGroupDetails{
				Parent: &managementgroups.CreateParentGroupInfo{
					ID: utils.String(parentID),
				},
			},
		},
		Type: utils.String("/providers/Microsoft.Management/managementGroups"),
		Name: &name,
	}

	log.Printf("[DEBUG] Invoking managementGroupClient")
	createManagementGroupFuture, err := client.CreateOrUpdate(ctx, name, properties, "no-cache")
	if err != nil {
		log.Printf("[DEBUG] Error creating Management Group %q: %+v", name, err)
		return fmt.Errorf("Error creating Management Group %q: %+v", name, err)
	}

	err = createManagementGroupFuture.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Management Group %q: %+v", name, err)
	}

	recurse := false

	resp, err := client.Get(ctx, name, "", &recurse, "", "no-cache")
	if err != nil {
		log.Printf("[DEBUG] Error retrieving Management Group %q: %+v", name, err)
		return fmt.Errorf("Error retrieving Management Group %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	for _, subscription := range subscriptionIds {
		data := subscription.(string)
		log.Printf("[DEBUG] Adding subscriptionId %q to management group %q", data, name)
		_, err = subscriptionsClient.Create(ctx, name, data, "no-cache")
		if err != nil {
			log.Printf("[DEBUG] Error assigning subscription %q to management group %q", data, name)
			return err
		}
	}

	return resourceManagementGroupRead(d, meta)
}

func resourceManagementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).managementGroupsClient
	ctx := meta.(*ArmClient).StopContext

	recurse := true
	resp, err := client.Get(ctx, d.Get("name").(string), "children", &recurse, "", "no-cache")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Management Group %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Management Group %+v", err)
	}

	subscriptionIds := []string{}

	if props := resp.Properties; props != nil {
		if children := props.Children; children != nil {
			for _, child := range *children {
				subscriptionID, err := parseSubscriptionID(*child.ID)
				if err != nil {
					log.Printf("%q", err)
					return fmt.Errorf("Unable to parse child subscription ID %+v", err)
				}
				log.Printf("[INFO] Reading subscription %q from management group %q", subscriptionID, d.Get("name").(string))
				subscriptionIds = append(subscriptionIds, subscriptionID)
			}
		}
	}

	d.Set("subscription_ids", subscriptionIds)

	return nil
}

func resourceManagementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	//before deleting a management group, return any subscriptions to the root management group

	client := meta.(*ArmClient).managementGroupsClient
	subscriptionsClient := meta.(*ArmClient).managementGroupsSubscriptionClient
	ctx := meta.(*ArmClient).StopContext
	armTenantID := meta.(*ArmClient).tenantId
	name := d.Get("name").(string)

	subscriptionIds := d.Get("subscription_ids").([]interface{})
	if subscriptionIds != nil {
		for _, subscription := range subscriptionIds {
			data := subscription.(string)
			log.Printf("[DEBUG] Adding subscriptionId %q to management group %q", data, armTenantID)
			_, err := subscriptionsClient.Create(ctx, armTenantID, data, "no-cache")
			if err != nil {
				log.Printf("[DEBUG] Error assigning subscription %q to management group %q", data, armTenantID)
				return err
			}
		}
	}

	resp, err := client.Delete(ctx, name, "no-cache")
	if err != nil {
		log.Printf("[DEBUG] Error deleting management group %q", name)
		return fmt.Errorf("Error deleting management group %q", name)
	}

	err = resp.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error deleting management group %q", name)
	}

	_, err = resp.Result(client)

	if err != nil {
		return fmt.Errorf("Error deleting management group %q", name)
	}

	return nil
}
func parseSubscriptionID(id string) (string, error) {
	components := strings.Split(id, "/")

	if len(components) == 0 {
		return "", fmt.Errorf("Subscription Id is empty or not formatted correctly: %s", id)
	}

	if len(components) != 3 {
		return "", fmt.Errorf("Subscription Id should have 2 segments, got %d: '%s'", len(components)-1, id)
	}

	return components[2], nil
}
