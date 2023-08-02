// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMachineLearningWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMachineLearningWorkspaceRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMachineLearningWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	d.Set("location", location.NormalizeNilable(resp.Model.Location))

	flattenedIdentity, err := flattenMachineLearningWorkspaceIdentity(resp.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Model.Tags)
}
