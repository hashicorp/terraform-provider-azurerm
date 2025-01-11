// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managementgroup

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managementgroups/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupSubscriptionAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagementGroupSubscriptionAssociationCreate,
		Read:   resourceManagementGroupSubscriptionAssociationRead,
		Delete: resourceManagementGroupSubscriptionAssociationDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ManagementGroupSubscriptionAssociationV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managementgroups.ParseSubscriptionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"management_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagementGroupID,
			},

			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubscriptionID,
			},
		},
	}
}

func resourceManagementGroupSubscriptionAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managementGroupId, err := parse.ManagementGroupID(d.Get("management_group_id").(string))
	if err != nil {
		return err
	}

	subscriptionId, err := commonids.ParseSubscriptionID(d.Get("subscription_id").(string))
	if err != nil {
		return err
	}

	id := managementgroups.NewSubscriptionID(managementGroupId.Name, subscriptionId.SubscriptionId)

	existing, err := client.Get(ctx, commonids.NewManagementGroupID(id.GroupId), managementgroups.GetOperationOptions{
		CacheControl: &managementGroupCacheControl,
		Expand:       pointer.To(managementgroups.ExpandChildren),
		Recurse:      pointer.FromBool(false),
	})
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("failed checking Management Group %q: %+v", id.GroupId, err)
		}
	}

	if model := existing.Model; model != nil {
		props := model.Properties
		if props == nil {
			return fmt.Errorf("could not read properties for Management Group %q to check if Subscription Association for %q already exists", id.GroupId, id.SubscriptionId)
		}

		if props.Children != nil {
			for _, v := range *props.Children {
				if v.Type != nil && *v.Type == managementgroups.ManagementGroupChildTypeSubscriptions && v.Name != nil && strings.EqualFold(*v.Name, id.SubscriptionId) {
					return tf.ImportAsExistsError("azurerm_management_group_subscription_association", id.ID())
				}
			}
		}
	}

	if _, err := client.SubscriptionsCreate(ctx, id, managementgroups.SubscriptionsCreateOperationOptions{}); err != nil {
		return fmt.Errorf("creating Management Group Subscription Association between %q and %q: %+v", managementGroupId.Name, subscriptionId, err)
	}

	d.SetId(id.ID())

	return resourceManagementGroupSubscriptionAssociationRead(d, meta)
}

func resourceManagementGroupSubscriptionAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	// There is no "read" function on the appropriate client so we need to check if the Subscription is in the Management Group subscription list
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managementgroups.ParseSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	managementGroup, err := client.Get(ctx, commonids.NewManagementGroupID(id.GroupId), managementgroups.GetOperationOptions{
		CacheControl: &managementGroupCacheControl,
		Expand:       pointer.To(managementgroups.ExpandChildren),
		Recurse:      pointer.FromBool(false),
	})
	if err != nil {
		return fmt.Errorf("reading Management Group %q for Subscription Associations: %+v", id.GroupId, err)
	}
	found := false
	if model := managementGroup.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.Children != nil {
				for _, v := range *props.Children {
					if v.Type != nil && *v.Type == managementgroups.ManagementGroupChildTypeSubscriptions && v.Name != nil && strings.EqualFold(*v.Name, id.SubscriptionId) {
						found = true
					}
				}
			}

			if !found {
				log.Printf("[INFO] Subscription %q not found in Management group %q, removing from state", id.SubscriptionId, id.GroupId)
				d.SetId("")
				return nil
			}

			managementGroupId := parse.NewManagementGroupId(id.GroupId)
			d.Set("management_group_id", managementGroupId.ID())
			subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)
			d.Set("subscription_id", subscriptionId.ID())
		}
	}

	return nil
}

func resourceManagementGroupSubscriptionAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementGroups.GroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managementgroups.ParseSubscriptionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SubscriptionsDelete(ctx, *id, managementgroups.SubscriptionsDeleteOperationOptions{})
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Management Group Subscription Association between Management Group %q and Subscription %q: %+v", id.GroupId, id.SubscriptionId, err)
		}
	}

	// It's a workaround to solve the replication delay issue: DELETE operation happens in one region, but it needs more time to sync the result to other regions.
	log.Printf("[DEBUG] Waiting for %s to be fully deleted..", d.Id())
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   subscriptionAssociationRefreshFunc(ctx, meta.(*clients.Client).ManagementGroups.GroupsClient, *id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be fully deleted: %+v", d.Id(), err)
	}

	return nil
}

func subscriptionAssociationRefreshFunc(ctx context.Context, client *managementgroups.ManagementGroupsClient, id managementgroups.SubscriptionId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		managementGroup, err := client.Get(ctx, commonids.NewManagementGroupID(id.GroupId), managementgroups.GetOperationOptions{
			CacheControl: &managementGroupCacheControl,
			Expand:       pointer.To(managementgroups.ExpandChildren),
			Recurse:      pointer.FromBool(false),
		})
		if err != nil {
			return nil, "", fmt.Errorf("reading Management Group %q for Subscription Associations: %+v", id.GroupId, err)
		}

		if model := managementGroup.Model; model != nil {
			if props := model.Properties; props != nil && props.Children != nil {
				for _, v := range *props.Children {
					if v.Type != nil && *v.Type == managementgroups.ManagementGroupChildTypeSubscriptions {
						if v.Name != nil && strings.EqualFold(*v.Name, id.SubscriptionId) {
							return managementGroup, "Exists", nil
						}
					}
				}
			}
		}

		return "NotFound", "NotFound", nil
	}
}
