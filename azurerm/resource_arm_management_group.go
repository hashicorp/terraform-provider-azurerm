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
				Required: false,
				Computed: false,
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

	createParentGroupInfo := managementgroups.CreateParentGroupInfo{
		ID: utils.String("/providers/Microsoft.Management/managementGroups/" + armTenantID),
	}
	details := managementgroups.CreateManagementGroupDetails{
		Parent: &createParentGroupInfo,
	}

	managementGroupProperties := managementgroups.CreateManagementGroupProperties{
		TenantID:    &armTenantID,
		DisplayName: &name,
		Details:     &details,
	}

	createManagementGroupRequest := managementgroups.CreateManagementGroupRequest{
		ID:   nil,
		Type: utils.String("/providers/Microsoft.Management/managementGroups"),
		Name: &name,
		CreateManagementGroupProperties: &managementGroupProperties,
	}

	log.Printf("[DEBUG] Invoking managementGroupClient")
	createManagementGroupFuture, err := client.CreateOrUpdate(ctx, name, createManagementGroupRequest, "no-cache")
	if err != nil {
		log.Printf("[DEBUG] Error invoking REST API call")
		return err
	}

	err = createManagementGroupFuture.WaitForCompletion(ctx, client.Client)

	_, err = createManagementGroupFuture.Result(client)
	if err != nil {
		log.Printf("[DEBUG] Error in API response")
		return err
	}

	recurse := false

	getResp, getErr := client.Get(ctx, name, "", &recurse, "", "no-cache")
	if getErr != nil {
		log.Printf("[DEBUG] Error retrieving management group details")
		return err
	}
	d.SetId(*getResp.ID)

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
	children := resp.Children
	if children != nil {
		for _, child := range *children {
			subscriptionID, err := parseSubscriptionID(*child.ID)
			if err != nil {
				log.Printf("%q", err)
			}
			log.Printf("[INFO] Reading subscription %q from management group %q", subscriptionID, d.Get("name").(string))
			subscriptionIds = append(subscriptionIds, subscriptionID)
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
