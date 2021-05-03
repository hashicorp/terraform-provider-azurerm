package managementgroup

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	subscriptionParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	subscriptionValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagementGroupSubscriptionAssociationID(id)
			return err
		}),

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
	groupsClient := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId, err := parse.ManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}

	subscriptionId, err := subscriptionParse.SubscriptionID(d.Get("subscription_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewManagementGroupSubscriptionAssociationID(managementGroupId.Name, subscriptionId.SubscriptionID)

	existing, err := groupsClient.Get(ctx, id.ManagementGroup, "children", utils.Bool(false), "", "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("failed checking Management Group %q: %+v", id.ManagementGroup, err)
		}
	}

	props := existing.Properties
	if props == nil {
		return fmt.Errorf("could not read properties for Management Group %q to check if Subscription Association for %q already exists", id.ManagementGroup, id.SubscriptionId)
	}

	if props.Children != nil {
		for _, v := range *props.Children {
			if v.Type == managementgroups.Type1Subscriptions && v.Name != nil && *v.Name == id.SubscriptionId {
				return tf.ImportAsExistsError("azurerm_management_group_subscription_association", id.ID())
			}
		}
	}

	if _, err := client.Create(ctx, id.ManagementGroup, id.SubscriptionId, ""); err != nil {
		return fmt.Errorf("creating Management Group Subscription Association between %q and %q: %+v", managementGroupId.Name, subscriptionId, err)
	}

	d.SetId(id.ID())

	return resourceManagementGroupSubscriptionAssociationRead(d, meta)
}

func resourceManagementGroupSubscriptionAssociationRead(d *schema.ResourceData, meta interface{}) error {
	// There is no "read" function on the appropriate client so we need to check if the Subscription is in the Management Group subscription list
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementGroupSubscriptionAssociationID(d.Id())
	if err != nil {
		return err
	}

	managementGroup, err := client.Get(ctx, id.ManagementGroup, "children", utils.Bool(false), "", "")
	if err != nil {
		return fmt.Errorf("reading Management Group %q for Subscription Associations: %+v", id.ManagementGroup, err)
	}
	found := false
	if props := managementGroup.Properties; props != nil {
		if props.Children == nil {
			return fmt.Errorf("could not read properties for Management Group %q", id.ManagementGroup)
		}

		for _, v := range *props.Children {
			if v.Type == managementgroups.Type1Subscriptions {
				if v.Name != nil && *v.Name == id.SubscriptionId {
					found = true
				}
			}
		}

		if !found {
			log.Printf("[INFO] Subscription %q not found in Management group %q, removing from state", id.SubscriptionId, id.ManagementGroup)
			d.SetId("")
			return nil
		}

		managementGroupId := parse.NewManagementGroupId(id.ManagementGroup)
		d.Set("management_group_id", managementGroupId.ID())
		subscriptionId := subscriptionParse.NewSubscriptionId(id.SubscriptionId)
		d.Set("subscription_id", subscriptionId.ID())
	}

	return nil
}

func resourceManagementGroupSubscriptionAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.SubscriptionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementGroupSubscriptionAssociationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ManagementGroup, id.SubscriptionId, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Management Group Subscription Association between Management Group %q and Subscription %q: %+v", id.ManagementGroup, id.SubscriptionId, err)
		}
	}

	return nil
}
