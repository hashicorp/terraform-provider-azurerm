package managementgroup

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var managementGroupCacheControl = "no-cache"

func resourceManagementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceManagementGroupCreateUpdate,
		Update: resourceManagementGroupCreateUpdate,
		Read:   resourceManagementGroupRead,
		Delete: resourceManagementGroupDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ManagementGroupID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				Deprecated:    "Deprecated in favour of `name`",
				ValidateFunc:  validate.ManagementGroupName,
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"group_id"},
				ValidateFunc:  validate.ManagementGroupName,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"parent_management_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ManagementGroupID,
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

func resourceManagementGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	subscriptionsClient := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	armTenantID := meta.(*clients.Client).Account.TenantId

	var groupName string
	if v := d.Get("group_id"); v != "" {
		groupName = v.(string)
	}
	if v := d.Get("name"); v != "" {
		groupName = v.(string)
	}

	if groupName == "" {
		groupName = uuid.New().String()
	}

	parentManagementGroupId := d.Get("parent_management_group_id").(string)
	if parentManagementGroupId == "" {
		parentManagementGroupId = fmt.Sprintf("/providers/Microsoft.Management/managementGroups/%s", armTenantID)
	}

	recurse := false
	if d.IsNewResource() {
		existing, err := client.Get(ctx, groupName, "children", &recurse, "", managementGroupCacheControl)
		if err != nil {
			// 403 is returned if group does not exist, bug tracked at: https://github.com/Azure/azure-rest-api-specs/issues/9549
			if !utils.ResponseWasNotFound(existing.Response) && !utils.ResponseWasForbidden(existing.Response) {
				return fmt.Errorf("unable to check for presence of existing Management Group %q: %s", groupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_management_group", *existing.ID)
		}
	}

	log.Printf("[INFO] Creating Management Group %q", groupName)

	properties := managementgroups.CreateManagementGroupRequest{
		Name: utils.String(groupName),
		CreateManagementGroupProperties: &managementgroups.CreateManagementGroupProperties{
			TenantID: utils.String(armTenantID),
			Details: &managementgroups.CreateManagementGroupDetails{
				Parent: &managementgroups.CreateParentGroupInfo{
					ID: utils.String(parentManagementGroupId),
				},
			},
		},
	}

	if v := d.Get("display_name"); v != "" {
		properties.CreateManagementGroupProperties.DisplayName = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, groupName, properties, managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("unable to create Management Group %q: %+v", groupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed when waiting for creation of Management Group %q: %+v", groupName, err)
	}

	// We have a potential race condition / consistency issue whereby the implicit role assignment for the SP may not be
	// completed before the read-back here or an eventually consistent read is creating a temporary 403 error.
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			"pending",
		},
		Target: []string{
			"succeeded",
		},
		Refresh:                   managementgroupCreateStateRefreshFunc(ctx, client, groupName),
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		ContinuousTargetOccurence: 5,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("failed waiting for read on Managementgroup %q", groupName)
	}

	resp, err := client.Get(ctx, groupName, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("unable to retrieve Management Group %q: %+v", groupName, err)
	}

	d.SetId(*resp.ID)

	subscriptionIds := expandManagementGroupSubscriptionIds(d.Get("subscription_ids").(*schema.Set))

	// first remove any which need to be removed
	if !d.IsNewResource() {
		log.Printf("[DEBUG] Determine which Subscriptions should be removed from Management Group %q", groupName)
		if props := resp.Properties; props != nil {
			subscriptionIdsToRemove, err2 := determineManagementGroupSubscriptionsIdsToRemove(props.Children, subscriptionIds)
			if err2 != nil {
				return fmt.Errorf("unable to determine which subscriptions should be removed from Management Group %q: %+v", groupName, err2)
			}

			for _, subscriptionId := range *subscriptionIdsToRemove {
				log.Printf("[DEBUG] De-associating Subscription ID %q from Management Group %q", subscriptionId, groupName)
				deleteResp, err2 := subscriptionsClient.Delete(ctx, groupName, subscriptionId, managementGroupCacheControl)
				if err2 != nil {
					if !response.WasNotFound(deleteResp.Response) {
						return fmt.Errorf("unable to de-associate Subscription %q from Management Group %q: %+v", subscriptionId, groupName, err2)
					}
				}
			}
		}
	}

	// then add the new ones
	log.Printf("[DEBUG] Preparing to assign Subscriptions to Management Group %q", groupName)
	for _, subscriptionId := range subscriptionIds {
		log.Printf("[DEBUG] Assigning Subscription ID %q to management group %q", subscriptionId, groupName)
		if _, err := subscriptionsClient.Create(ctx, groupName, subscriptionId, managementGroupCacheControl); err != nil {
			return fmt.Errorf("[DEBUG] Error assigning Subscription ID %q to Management Group %q: %+v", subscriptionId, groupName, err)
		}
	}

	return resourceManagementGroupRead(d, meta)
}

func resourceManagementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name, "children", utils.Bool(true), "", managementGroupCacheControl)
	if err != nil {
		if utils.ResponseWasForbidden(resp.Response) || utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Management Group %q doesn't exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("unable to read Management Group %q: %+v", d.Id(), err)
	}

	d.Set("name", id.Name)
	d.Set("group_id", id.Name)

	if props := resp.Properties; props != nil {
		d.Set("display_name", props.DisplayName)

		subscriptionIds, err := flattenManagementGroupSubscriptionIds(props.Children)
		if err != nil {
			return fmt.Errorf("unable to flatten `subscription_ids`: %+v", err)
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

func resourceManagementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	subscriptionsClient := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementGroupID(d.Id())
	if err != nil {
		return err
	}

	recurse := true
	group, err := client.Get(ctx, id.Name, "children", &recurse, "", managementGroupCacheControl)
	if err != nil {
		if utils.ResponseWasNotFound(group.Response) || utils.ResponseWasForbidden(group.Response) {
			log.Printf("[DEBUG] Management Group %q doesn't exist in Azure - nothing to do!", id.Name)
			return nil
		}

		return fmt.Errorf("unable to retrieve Management Group %q: %+v", id.Name, err)
	}

	// before deleting a management group, return any subscriptions to the root management group
	if props := group.Properties; props != nil {
		if children := props.Children; children != nil {
			for _, v := range *children {
				if v.ID == nil {
					continue
				}

				subscriptionId, err := parseManagementGroupSubscriptionID(*v.ID)
				if err != nil {
					return fmt.Errorf("unable to parse child Subscription ID %+v", err)
				}
				log.Printf("[DEBUG] De-associating Subscription %q from Management Group %q..", subscriptionId, id.Name)
				// NOTE: whilst this says `Delete` it's actually `Deassociate` - which is /really/ helpful
				deleteResp, err2 := subscriptionsClient.Delete(ctx, id.Name, subscriptionId.subscriptionId, managementGroupCacheControl)
				if err2 != nil {
					if !response.WasNotFound(deleteResp.Response) {
						return fmt.Errorf("unable to de-associate Subscription %q from Management Group %q: %+v", subscriptionId.subscriptionId, id.Name, err2)
					}
				}
			}
		}
	}

	resp, err := client.Delete(ctx, id.Name, managementGroupCacheControl)
	if err != nil {
		return fmt.Errorf("unable to delete Management Group %q: %+v", id.Name, err)
	}

	err = resp.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("failed when waiting for the deletion of Management Group %q: %+v", id.Name, err)
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

func flattenManagementGroupSubscriptionIds(input *[]managementgroups.ChildInfo) (*schema.Set, error) {
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
			return nil, fmt.Errorf("unable to parse child Subscription ID %+v", err)
		}

		if id != nil {
			subscriptionIds.Add(id.subscriptionId)
		}
	}

	return subscriptionIds, nil
}

type subscriptionId struct {
	subscriptionId string
}

func parseManagementGroupSubscriptionID(input string) (*subscriptionId, error) {
	// this is either:
	// /subscriptions/00000000-0000-0000-0000-000000000000

	// we skip out the managementGroup ID's
	if strings.HasPrefix(input, "/providers/Microsoft.Management/managementGroups/") {
		return nil, nil
	}

	components := strings.Split(input, "/")

	if len(components) == 0 {
		return nil, fmt.Errorf("subscription Id is empty or not formatted correctly: %s", input)
	}

	if len(components) != 3 {
		return nil, fmt.Errorf("subscription Id should have 2 segments, got %d: %q", len(components)-1, input)
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
			return nil, fmt.Errorf("unable to parse Subscription ID %q: %+v", *v.ID, err)
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

func managementgroupCreateStateRefreshFunc(ctx context.Context, client *managementgroups.Client, groupName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, groupName, "children", utils.Bool(true), "", managementGroupCacheControl)
		if err != nil {
			if utils.ResponseWasForbidden(resp.Response) {
				return resp, "pending", nil
			}
			return resp, "failed", err
		}

		return resp, "succeeded", nil
	}
}
