package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var managementGroupCacheControl = "no-cache"

func resourceArmManagementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagementGroupCreateUpdate,
		Update: resourceArmManagementGroupCreateUpdate,
		Read:   resourceArmManagementGroupRead,
		Delete: resourceArmManagementGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"parent_management_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"subscription_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceArmManagementGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	subscriptionsClient := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	armTenantID := meta.(*clients.Client).Account.TenantId

	groupId := d.Get("group_id").(string)
	if groupId == "" {
		groupId = uuid.New().String()
	}

	parentManagementGroupId := d.Get("parent_management_group_id").(string)
	if parentManagementGroupId == "" {
		parentManagementGroupId = fmt.Sprintf("/providers/Microsoft.Management/managementGroups/%s", armTenantID)
	}

	recurse := false
	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, groupId, "children", &recurse, "", managementGroupCacheControl)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Management Group %q: %s", groupId, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_management_group", *existing.ID)
		}
	}

	log.Printf("[INFO] Creating Management Group %q", groupId)

	properties := managementgroups.CreateManagementGroupRequest{
		Name: utils.String(groupId),
		CreateManagementGroupProperties: &managementgroups.CreateManagementGroupProperties{
			TenantID: utils.String(armTenantID),
			Details: &managementgroups.CreateManagementGroupDetails{
				Parent: &managementgroups.CreateParentGroupInfo{
					ID: utils.String(parentManagementGroupId),
				},
			},
		},
	}

	if v, ok := d.GetOk("display_name"); ok {
		properties.CreateManagementGroupProperties.DisplayName = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, groupId, properties, managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("Error creating Management Group %q: %+v", groupId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Management Group %q: %+v", groupId, err)
	}

	resp, err := client.Get(ctx, groupId, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("Error retrieving Management Group %q: %+v", groupId, err)
	}

	d.SetId(*resp.ID)

	subscriptionIds := expandManagementGroupSubscriptionIds(d.Get("subscription_ids").(*schema.Set))

	// first remove any which need to be removed
	if !d.IsNewResource() {
		log.Printf("[DEBUG] Determine which Subscriptions should be removed from Management Group %q", groupId)
		if props := resp.Properties; props != nil {
			subscriptionIdsToRemove, err2 := determineManagementGroupSubscriptionsIdsToRemove(props.Children, subscriptionIds)
			if err2 != nil {
				return fmt.Errorf("Error determining which subscriptions should be removed from Management Group %q: %+v", groupId, err2)
			}

			for _, subscriptionId := range *subscriptionIdsToRemove {
				log.Printf("[DEBUG] De-associating Subscription ID %q from Management Group %q", subscriptionId, groupId)
				deleteResp, err2 := subscriptionsClient.Delete(ctx, groupId, subscriptionId, managementGroupCacheControl)
				if err2 != nil {
					if !response.WasNotFound(deleteResp.Response) {
						return fmt.Errorf("Error de-associating Subscription %q from Management Group %q: %+v", subscriptionId, groupId, err2)
					}
				}
			}
		}
	}

	// then add the new ones
	log.Printf("[DEBUG] Preparing to assign Subscriptions to Management Group %q", groupId)
	for _, subscriptionId := range subscriptionIds {
		log.Printf("[DEBUG] Assigning Subscription ID %q to management group %q", subscriptionId, groupId)
		_, err = subscriptionsClient.Create(ctx, groupId, subscriptionId, managementGroupCacheControl)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error assigning Subscription ID %q to Management Group %q: %+v", subscriptionId, groupId, err)
		}
	}

	return resourceArmManagementGroupRead(d, meta)
}

func resourceArmManagementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseManagementGroupId(d.Id())
	if err != nil {
		return err
	}

	recurse := true
	resp, err := client.Get(ctx, id.groupId, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Management Group %q doesn't exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Management Group %q: %+v", d.Id(), err)
	}

	d.Set("group_id", id.groupId)

	if props := resp.Properties; props != nil {
		d.Set("display_name", props.DisplayName)

		subscriptionIds, err := flattenArmManagementGroupSubscriptionIds(props.Children)
		if err != nil {
			return fmt.Errorf("Error flattening `subscription_ids`: %+v", err)
		}
		d.Set("subscription_ids", subscriptionIds)

		parentId := ""
		if details := props.Details; details != nil {
			if parent := details.Parent; parent != nil {
				if pid := parent.ID; pid != nil {
					parentId = *pid
				}
			}
		}
		d.Set("parent_management_group_id", parentId)
	}

	return nil
}

func resourceArmManagementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	subscriptionsClient := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseManagementGroupId(d.Id())
	if err != nil {
		return err
	}

	recurse := true
	group, err := client.Get(ctx, id.groupId, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		if utils.ResponseWasNotFound(group.Response) {
			log.Printf("[DEBUG] Management Group %q doesn't exist in Azure - nothing to do!", id.groupId)
			return nil
		}

		return fmt.Errorf("Error retrieving Management Group %q: %+v", id.groupId, err)
	}

	// before deleting a management group, return any subscriptions to the root management group
	if props := group.Properties; props != nil {
		if children := props.Children; children != nil {
			for _, v := range *children {
				if v.ID == nil {
					continue
				}

				subscriptionId := *v.ID
				log.Printf("[DEBUG] De-associating Subscription %q from Management Group %q..", subscriptionId, id.groupId)
				// NOTE: whilst this says `Delete` it's actually `Deassociate` - which is /really/ helpful
				deleteResp, err2 := subscriptionsClient.Delete(ctx, id.groupId, subscriptionId, managementGroupCacheControl)
				if err2 != nil {
					if !response.WasNotFound(deleteResp.Response) {
						return fmt.Errorf("Error de-associating Subscription %q from Management Group %q: %+v", subscriptionId, id.groupId, err2)
					}
				}
			}
		}
	}

	resp, err := client.Delete(ctx, id.groupId, managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("Error deleting Management Group %q: %+v", id.groupId, err)
	}

	err = resp.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of Management Group %q: %+v", id.groupId, err)
	}

	return nil
}

func expandManagementGroupSubscriptionIds(input *schema.Set) []string {
	output := make([]string, 0)

	if input != nil {
		for _, v := range input.List() {
			output = append(output, v.(string))
		}
	}

	return output
}

func flattenArmManagementGroupSubscriptionIds(input *[]managementgroups.ChildInfo) (*schema.Set, error) {
	subscriptionIds := &schema.Set{F: schema.HashString}
	if input == nil {
		return subscriptionIds, nil
	}

	for _, child := range *input {
		if child.ID == nil {
			continue
		}

		id, err := parseManagementGroupSubscriptionID(*child.ID)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse child Subscription ID %+v", err)
		}

		if id != nil {
			subscriptionIds.Add(id.subscriptionId)
		}
	}

	return subscriptionIds, nil
}

type managementGroupId struct {
	groupId string
}

type subscriptionId struct {
	subscriptionId string
}

func parseManagementGroupId(input string) (*managementGroupId, error) {
	// /providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000
	segments := strings.Split(input, "/")
	if len(segments) != 5 {
		return nil, fmt.Errorf("Expected there to be 5 segments but got %d", len(segments))
	}

	id := managementGroupId{
		groupId: segments[4],
	}
	return &id, nil
}

func parseManagementGroupSubscriptionID(input string) (*subscriptionId, error) {
	// this is either:
	// /subscriptions/00000000-0000-0000-0000-000000000000
	// /providers/Microsoft.Management/managementGroups/e4115b99-6be7-4153-a73f-5ff5e778ce28

	// we skip out the managementGroup ID's
	if strings.HasPrefix(input, "/providers/Microsoft.Management/managementGroups/") {
		return nil, nil
	}

	components := strings.Split(input, "/")

	if len(components) == 0 {
		return nil, fmt.Errorf("Subscription Id is empty or not formatted correctly: %s", input)
	}

	if len(components) != 3 {
		return nil, fmt.Errorf("Subscription Id should have 2 segments, got %d: %q", len(components)-1, input)
	}

	id := subscriptionId{
		subscriptionId: components[2],
	}
	return &id, nil
}

func determineManagementGroupSubscriptionsIdsToRemove(existing *[]managementgroups.ChildInfo, updated []string) (*[]string, error) {
	subscriptionIdsToRemove := make([]string, 0)
	if existing == nil {
		return &subscriptionIdsToRemove, nil
	}

	for _, v := range *existing {
		if v.ID == nil {
			continue
		}

		id, err := parseManagementGroupSubscriptionID(*v.ID)
		if err != nil {
			return nil, fmt.Errorf("Error parsing Subscription ID %q: %+v", *v.ID, err)
		}

		// not a Subscription - so let's skip it
		if id == nil {
			continue
		}

		found := false
		for _, subId := range updated {
			if id.subscriptionId == subId {
				found = true
				break
			}
		}

		if !found {
			subscriptionIdsToRemove = append(subscriptionIdsToRemove, id.subscriptionId)
		}
	}

	return &subscriptionIdsToRemove, nil
}
