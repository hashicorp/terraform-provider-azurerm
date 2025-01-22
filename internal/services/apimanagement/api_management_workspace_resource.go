// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementWorkspaceCreateUpdate,
		Read:   resourceApiManagementWorkspaceRead,
		Update: resourceApiManagementWorkspaceCreateUpdate,
		Delete: resourceApiManagementWorkspaceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workspace.ParseWorkspaceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(45 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": commonschema.ResourceGroupName(),
			"workspace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceApiManagementWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.WorkspaceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for API Management Workspace creation.")

	id := workspace.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("service_name").string, d.Get("workspace_id").(string))

	//...
	//TODO ctx

	properties := workspace.WorkspaceContract{
		//TODO
	}
	//...

	if _, err := client.CreateOrUpdate(ctx, id, properties, workspace.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	return resourceApiManagementWorkspaceDelete(d, meta)
}

func resourceApiManagementWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.WorkspaceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspace.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("workspace_id", id.WorkspaceId)
	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	return nil
}

func resourceApiManagementWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
