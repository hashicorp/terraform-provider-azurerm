// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/groupuser"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGroupUser() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGroupUserCreate,
		Read:   resourceApiManagementGroupUserRead,
		Delete: resourceApiManagementGroupUserDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := groupuser.ParseGroupUserID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"user_id": schemaz.SchemaApiManagementChildName(),

			"group_name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),
		},
	}
}

func resourceApiManagementGroupUserCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := groupuser.NewGroupUserID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("group_name").(string), d.Get("user_id").(string))

	exists, err := client.CheckEntityExists(ctx, id)
	if err != nil {
		if !response.WasNotFound(exists.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(exists.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_api_management_group_user", id.ID())
	}

	if _, err := client.Create(ctx, id); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGroupUserRead(d, meta)
}

func resourceApiManagementGroupUserRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := groupuser.ParseGroupUserID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.CheckEntityExists(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("group_name", id.GroupId)
	d.Set("user_id", id.UserId)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	return nil
}

func resourceApiManagementGroupUserDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := groupuser.ParseGroupUserID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
