// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/webhooks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceContainerRegistryWebhook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryWebhookCreate,
		Read:   resourceContainerRegistryWebhookRead,
		Update: resourceContainerRegistryWebhookUpdate,
		Delete: resourceContainerRegistryWebhookDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webhooks.ParseWebHookID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.RegistryWebhookV0ToV1{},
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
				ValidateFunc: validate.ContainerRegistryWebhookName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"registry_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ContainerRegistryName,
			},

			"service_uri": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ContainerRegistryWebhookServiceUri,
			},

			"custom_headers": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"status": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  webhooks.WebhookStatusEnabled,
				ValidateFunc: validation.StringInSlice([]string{
					string(webhooks.WebhookStatusDisabled),
					string(webhooks.WebhookStatusEnabled),
				}, false),
			},

			"scope": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "",
			},

			"actions": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(webhooks.WebhookActionChartDelete),
						string(webhooks.WebhookActionChartPush),
						string(webhooks.WebhookActionDelete),
						string(webhooks.WebhookActionPush),
						string(webhooks.WebhookActionQuarantine),
					}, false),
				},
			},

			"location": commonschema.Location(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceContainerRegistryWebhookCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.WebHooks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Container Registry Webhook creation.")

	id := webhooks.NewWebHookID(subscriptionId, d.Get("resource_group_name").(string), d.Get("registry_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_container_registry_webhook", id.ID())
		}
	}

	webhook := webhooks.WebhookCreateParameters{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: expandWebhookPropertiesCreateParameters(d),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, id, webhook); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryWebhookRead(d, meta)
}

func resourceContainerRegistryWebhookUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.WebHooks
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Container Registry Webhook update.")

	id, err := webhooks.ParseWebHookID(d.Id())
	if err != nil {
		return err
	}

	webhook := webhooks.WebhookUpdateParameters{
		Properties: expandWebhookPropertiesUpdateParameters(d),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.UpdateThenPoll(ctx, *id, webhook); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceContainerRegistryWebhookRead(d, meta)
}

func resourceContainerRegistryWebhookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.WebHooks
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webhooks.ParseWebHookID(d.Id())
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

	callbackConfig, err := client.GetCallbackConfig(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Callback Config for %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("registry_name", id.RegistryName)
	d.Set("name", id.WebHookName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			status := ""
			if v := props.Status; v != nil {
				status = string(*v)
			}
			d.Set("status", status)

			scope := ""
			if v := props.Scope; v != nil {
				scope = *v
			}
			d.Set("scope", scope)

			webhookActions := make([]string, len(props.Actions))
			for i, action := range props.Actions {
				webhookActions[i] = string(action)
			}
			d.Set("actions", webhookActions)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	if callbackModel := callbackConfig.Model; callbackModel != nil {
		if props := callbackModel; props != nil {
			d.Set("service_uri", props.ServiceUri)

			customHeaders := make(map[string]string)
			if props.CustomHeaders != nil {
				for k, v := range *props.CustomHeaders {
					customHeaders[k] = v
				}
			}
			d.Set("custom_headers", customHeaders)
		}
	}
	return nil
}

func resourceContainerRegistryWebhookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2021_08_01_preview.WebHooks
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webhooks.ParseWebHookID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandWebhookPropertiesCreateParameters(d *pluginsdk.ResourceData) *webhooks.WebhookPropertiesCreateParameters {
	customHeaders := make(map[string]string)
	for k, v := range d.Get("custom_headers").(map[string]interface{}) {
		customHeaders[k] = v.(string)
	}

	webhookProperties := webhooks.WebhookPropertiesCreateParameters{
		ServiceUri:    d.Get("service_uri").(string),
		CustomHeaders: &customHeaders,
		Actions:       expandWebhookActions(d),
		Scope:         pointer.To(d.Get("scope").(string)),
		Status:        pointer.To(webhooks.WebhookStatus(d.Get("status").(string))),
	}

	return &webhookProperties
}

func expandWebhookPropertiesUpdateParameters(d *pluginsdk.ResourceData) *webhooks.WebhookPropertiesUpdateParameters {
	customHeaders := make(map[string]string)
	for k, v := range d.Get("custom_headers").(map[string]interface{}) {
		customHeaders[k] = v.(string)
	}

	webhookProperties := webhooks.WebhookPropertiesUpdateParameters{
		ServiceUri:    pointer.To(d.Get("service_uri").(string)),
		CustomHeaders: &customHeaders,
		Actions:       pointer.To(expandWebhookActions(d)),
		Scope:         pointer.To(d.Get("scope").(string)),
		Status:        pointer.To(webhooks.WebhookStatus(d.Get("status").(string))),
	}

	return &webhookProperties
}

func expandWebhookActions(d *pluginsdk.ResourceData) []webhooks.WebhookAction {
	actions := make([]webhooks.WebhookAction, 0)
	for _, action := range d.Get("actions").(*pluginsdk.Set).List() {
		actions = append(actions, webhooks.WebhookAction(action.(string)))
	}

	return actions
}
