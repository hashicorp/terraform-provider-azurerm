// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_subscription_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/management/2020-05-01/managementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementGroupSubscriptionAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
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
		Recurse:      pointer.To(false),
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
					if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError("azurerm_management_group_subscription_association", id.ID())
					}
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
