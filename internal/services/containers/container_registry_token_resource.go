// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/scopemaps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceContainerRegistryToken() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryTokenCreate,
		Read:   resourceContainerRegistryTokenRead,
		Update: resourceContainerRegistryTokenUpdate,
		Delete: resourceContainerRegistryTokenDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := tokens.ParseTokenID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ContainerRegistryTokenName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"container_registry_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ContainerRegistryName,
			},

			"scope_map_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: scopemaps.ValidateScopeMapID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceContainerRegistryTokenCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Tokens
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := tokens.NewTokenID(subscriptionId, d.Get("resource_group_name").(string), d.Get("container_registry_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_token", id.ID())
		}
	}

	scopeMapID := d.Get("scope_map_id").(string)
	enabled := d.Get("enabled").(bool)
	status := tokens.TokenStatusEnabled

	if !enabled {
		status = tokens.TokenStatusDisabled
	}

	parameters := tokens.Token{
		Properties: &tokens.TokenProperties{
			ScopeMapId: utils.String(scopeMapID),
			Status:     &status,
		},
	}

	if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryTokenRead(d, meta)
}

func resourceContainerRegistryTokenUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Tokens
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Container Registry token update.")
	id, err := tokens.ParseTokenID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	scopeMapID := d.Get("scope_map_id").(string)
	enabled := d.Get("enabled").(bool)
	status := tokens.TokenStatusEnabled

	if !enabled {
		status = tokens.TokenStatusDisabled
	}

	parameters := tokens.TokenUpdateParameters{
		Properties: &tokens.TokenUpdateProperties{
			ScopeMapId: utils.String(scopeMapID),
			Status:     &status,
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryTokenRead(d, meta)
}

func resourceContainerRegistryTokenRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Tokens
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tokens.ParseTokenID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.TokenName)
	d.Set("container_registry_name", id.RegistryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			status := true
			if v := props.Status; v != nil && *v == tokens.TokenStatusDisabled {
				status = false
			}
			d.Set("enabled", status)

			scopeMapId := ""
			if v := props.ScopeMapId; v != nil {
				scopeMapId = *v
			}
			d.Set("scope_map_id", scopeMapId)
		}
	}

	return nil
}

func resourceContainerRegistryTokenDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.Tokens
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tokens.ParseTokenID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
