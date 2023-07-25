// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	validate2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceContainerRegistryAgentPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryAgentPoolCreate,
		Read:   resourceContainerRegistryAgentPoolRead,
		Update: resourceContainerRegistryAgentPoolUpdate,
		Delete: resourceContainerRegistryAgentPoolDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := agentpools.ParseAgentPoolID(id)
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
				ValidateFunc: validation.StringLenBetween(3, 20),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"container_registry_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate2.ContainerRegistryName,
			},

			"instance_count": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  1,
			},

			"tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "S1",
				ValidateFunc: validation.StringInSlice([]string{
					"S1",
					"S2",
					"S3",
					"I6",
				}, false),
			},

			"virtual_network_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceContainerRegistryAgentPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2019_06_01_preview.AgentPools
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Container Registry Agent Pool creation.")

	id := agentpools.NewAgentPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("container_registry_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_agent_pool", id.ID())
		}
	}

	parameters := agentpools.AgentPool{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &agentpools.AgentPoolProperties{
			// @favoretti: Only Linux is supported
			Os:    pointer.To(agentpools.OSLinux),
			Count: utils.Int64(int64(d.Get("instance_count").(int))),
			Tier:  pointer.To(d.Get("tier").(string)),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("virtual_network_subnet_id"); ok {
		parameters.Properties.VirtualNetworkSubnetResourceId = pointer.To(v.(string))
	}

	if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)

	}

	d.SetId(id.ID())

	return resourceContainerRegistryAgentPoolRead(d, meta)
}

func resourceContainerRegistryAgentPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2019_06_01_preview.AgentPools
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Container Registry Agent Pool creation.")

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	parameters := agentpools.AgentPoolUpdateParameters{
		Properties: &agentpools.AgentPoolPropertiesUpdateParameters{
			Count: utils.Int64(int64(d.Get("instance_count").(int))),
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryAgentPoolRead(d, meta)
}

func resourceContainerRegistryAgentPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2019_06_01_preview.AgentPools
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
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

	d.Set("name", id.AgentPoolName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("container_registry_name", id.RegistryName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			count := int64(0)
			if v := props.Count; v != nil {
				count = *v
			}
			d.Set("instance_count", count)

			tier := ""
			if v := props.Tier; v != nil {
				tier = *v
			}
			d.Set("tier", tier)

			virtualNetworkSubnetId := ""
			if v := props.VirtualNetworkSubnetResourceId; v != nil {
				virtualNetworkSubnetId = *v
			}
			d.Set("virtual_network_subnet_id", virtualNetworkSubnetId)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceContainerRegistryAgentPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2019_06_01_preview.AgentPools
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
