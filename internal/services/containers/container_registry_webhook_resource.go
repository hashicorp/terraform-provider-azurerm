package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2021-08-01-preview/containerregistry"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceContainerRegistryWebhook() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceContainerRegistryWebhookCreate,
		Read:   resourceContainerRegistryWebhookRead,
		Update: resourceContainerRegistryWebhookUpdate,
		Delete: resourceContainerRegistryWebhookDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WebhookID(id)
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
				ValidateFunc: validate.ContainerRegistryWebhookName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				Default:  containerregistry.WebhookStatusEnabled,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerregistry.WebhookStatusDisabled),
					string(containerregistry.WebhookStatusEnabled),
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
						string(containerregistry.WebhookActionChartDelete),
						string(containerregistry.WebhookActionChartPush),
						string(containerregistry.WebhookActionDelete),
						string(containerregistry.WebhookActionPush),
						string(containerregistry.WebhookActionQuarantine),
					}, false),
				},
			},

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceContainerRegistryWebhookCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.WebhooksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for  Container Registry Webhook creation.")

	id := parse.NewWebhookID(subscriptionId, d.Get("resource_group_name").(string), d.Get("registry_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_container_registry_webhook", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	webhook := containerregistry.WebhookCreateParameters{
		Location:                          &location,
		WebhookPropertiesCreateParameters: expandWebhookPropertiesCreateParameters(d),
		Tags:                              tags.Expand(t),
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.RegistryName, id.Name, webhook)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceContainerRegistryWebhookRead(d, meta)
}

func resourceContainerRegistryWebhookUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.WebhooksClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for  Container Registry Webhook update.")

	id, err := parse.WebhookID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	webhook := containerregistry.WebhookUpdateParameters{
		WebhookPropertiesUpdateParameters: expandWebhookPropertiesUpdateParameters(d),
		Tags:                              tags.Expand(t),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.RegistryName, id.Name, webhook)
	if err != nil {
		return fmt.Errorf("updating Container Registry Webhook %q (Resource Group %q, Registry %q): %+v", id.Name, id.ResourceGroup, id.RegistryName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Container Registry Webhook %q (Resource Group %q, Registry %q): %+v", id.Name, id.ResourceGroup, id.RegistryName, err)
	}

	return resourceContainerRegistryWebhookRead(d, meta)
}

func resourceContainerRegistryWebhookRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.WebhooksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebhookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Registry Webhook %q was not found in Resource Group %q for Registry %q", id.Name, id.ResourceGroup, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure Container Registry Webhook %q (Resource Group %q, Registry %q): %+v", id.Name, id.ResourceGroup, id.RegistryName, err)
	}

	callbackConfig, err := client.GetCallbackConfig(ctx, id.ResourceGroup, id.RegistryName, id.Name)
	if err != nil {
		return fmt.Errorf("making Read request on Azure Container Registry Webhook Callback Config %q (Resource Group %q, Registry %q): %+v", id.Name, id.ResourceGroup, id.RegistryName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("registry_name", id.RegistryName)
	d.Set("name", id.Name)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("service_uri", callbackConfig.ServiceURI)

	if callbackConfig.CustomHeaders != nil {
		customHeaders := make(map[string]string)
		for k, v := range callbackConfig.CustomHeaders {
			customHeaders[k] = *v
		}
		d.Set("custom_headers", customHeaders)
	}

	if webhookProps := resp.WebhookProperties; webhookProps != nil {
		if webhookProps.Status != "" {
			d.Set("status", string(webhookProps.Status))
		}

		if webhookProps.Scope != nil {
			d.Set("scope", webhookProps.Scope)
		}

		webhookActions := make([]string, len(*webhookProps.Actions))
		for i, action := range *webhookProps.Actions {
			webhookActions[i] = string(action)
		}
		d.Set("actions", webhookActions)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceContainerRegistryWebhookDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.WebhooksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebhookID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.RegistryName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Webhook %q (Container Registry %q / Resource Group %q): %+v", id.Name, id.RegistryName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Webhook %q (Container Registry %q / Resource Group %q): %+v", id.Name, id.RegistryName, id.ResourceGroup, err)
	}
	return nil
}

func expandWebhookPropertiesCreateParameters(d *pluginsdk.ResourceData) *containerregistry.WebhookPropertiesCreateParameters {
	serviceUri := d.Get("service_uri").(string)
	scope := d.Get("scope").(string)

	customHeaders := make(map[string]*string)
	for k, v := range d.Get("custom_headers").(map[string]interface{}) {
		customHeaders[k] = utils.String(v.(string))
	}

	actions := expandWebhookActions(d)

	webhookProperties := containerregistry.WebhookPropertiesCreateParameters{
		ServiceURI:    &serviceUri,
		CustomHeaders: customHeaders,
		Actions:       actions,
		Scope:         &scope,
	}

	webhookProperties.Status = containerregistry.WebhookStatus(d.Get("status").(string))

	return &webhookProperties
}

func expandWebhookPropertiesUpdateParameters(d *pluginsdk.ResourceData) *containerregistry.WebhookPropertiesUpdateParameters {
	serviceUri := d.Get("service_uri").(string)
	scope := d.Get("scope").(string)

	customHeaders := make(map[string]*string)
	for k, v := range d.Get("custom_headers").(map[string]interface{}) {
		customHeaders[k] = utils.String(v.(string))
	}

	webhookProperties := containerregistry.WebhookPropertiesUpdateParameters{
		ServiceURI:    &serviceUri,
		CustomHeaders: customHeaders,
		Actions:       expandWebhookActions(d),
		Scope:         &scope,
		Status:        containerregistry.WebhookStatus(d.Get("status").(string)),
	}

	return &webhookProperties
}

func expandWebhookActions(d *pluginsdk.ResourceData) *[]containerregistry.WebhookAction {
	actions := make([]containerregistry.WebhookAction, 0)
	for _, action := range d.Get("actions").(*pluginsdk.Set).List() {
		actions = append(actions, containerregistry.WebhookAction(action.(string)))
	}

	return &actions
}
