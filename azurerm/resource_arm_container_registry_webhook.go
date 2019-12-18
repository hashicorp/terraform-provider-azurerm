package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2018-09-01/containerregistry"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerRegistryWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerRegistryWebhookCreate,
		Read:   resourceArmContainerRegistryWebhookRead,
		Update: resourceArmContainerRegistryWebhookUpdate,
		Delete: resourceArmContainerRegistryWebhookDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMContainerRegistryWebhookName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"registry_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMContainerRegistryName,
			},

			"service_uri": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAzureRMContainerRegistryWebhookServiceUri,
			},

			"custom_headers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  containerregistry.WebhookStatusEnabled,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerregistry.WebhookStatusDisabled),
					string(containerregistry.WebhookStatusEnabled),
				}, false),
			},

			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"actions": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(containerregistry.ChartDelete),
						string(containerregistry.ChartPush),
						string(containerregistry.Delete),
						string(containerregistry.Push),
						string(containerregistry.Quarantine),
					}, false),
				},
			},

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmContainerRegistryWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Containers.WebhooksClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Container Registry Webhook creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	registryName := d.Get("registry_name").(string)
	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, registryName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Container Registry Webhook %q (Resource Group %q, Registry %q): %s", name, resourceGroup, registryName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_container_registry_webhook", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	webhook := containerregistry.WebhookCreateParameters{
		Location:                          &location,
		WebhookPropertiesCreateParameters: expandWebhookPropertiesCreateParameters(d),
		Tags:                              tags.Expand(t),
	}

	future, err := client.Create(ctx, resourceGroup, registryName, name, webhook)
	if err != nil {
		return fmt.Errorf("Error creating Container Registry Webhook %q (Resource Group %q, Registry %q): %+v", name, resourceGroup, registryName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Container Registry %q (Resource Group %q, Registry %q): %+v", name, resourceGroup, registryName, err)
	}

	read, err := client.Get(ctx, resourceGroup, registryName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Container Registry %q (Resource Group %q, Registry %q): %+v", name, resourceGroup, registryName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Container Registry %q (resource group %q, Registry %q) ID", name, resourceGroup, registryName)
	}

	d.SetId(*read.ID)

	return resourceArmContainerRegistryWebhookRead(d, meta)
}

func resourceArmContainerRegistryWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Containers.WebhooksClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Container Registry Webhook update.")

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	registryName := id.Path["registries"]
	name := id.Path["webhooks"]

	t := d.Get("tags").(map[string]interface{})

	webhook := containerregistry.WebhookUpdateParameters{
		WebhookPropertiesUpdateParameters: expandWebhookPropertiesUpdateParameters(d),
		Tags:                              tags.Expand(t),
	}

	future, err := client.Update(ctx, resourceGroup, registryName, name, webhook)
	if err != nil {
		return fmt.Errorf("Error updating Container Registry Webhook %q (Resource Group %q, Registry %q): %+v", name, resourceGroup, registryName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Container Registry Webhook %q (Resource Group %q, Registry %q): %+v", name, resourceGroup, registryName, err)
	}

	return resourceArmContainerRegistryWebhookRead(d, meta)
}

func resourceArmContainerRegistryWebhookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Containers.WebhooksClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	registryName := id.Path["registries"]
	name := id.Path["webhooks"]

	resp, err := client.Get(ctx, resourceGroup, registryName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Registry Webhook %q was not found in Resource Group %q for Registry %q", name, resourceGroup, registryName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Container Registry Webhook %q (Resource Group %q, Registry %q): %+v", name, resourceGroup, registryName, err)
	}

	callbackConfig, err := client.GetCallbackConfig(ctx, resourceGroup, registryName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Container Registry Webhook Callback Config %q (Resource Group %q, Registry %q): %+v", name, resourceGroup, registryName, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("registry_name", registryName)
	d.Set("name", name)

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

func resourceArmContainerRegistryWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Containers.WebhooksClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	registryName := id.Path["registries"]
	name := id.Path["webhooks"]

	future, err := client.Delete(ctx, resourceGroup, registryName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry Webhook '%s': %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing Azure ARM delete request of Container Registry Webhook '%s': %+v", name, err)
	}

	return nil
}

func validateAzureRMContainerRegistryWebhookName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9]{5,50}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"alpha numeric characters only are allowed and between 5 and 50 characters in %q: %q", k, value))
	}

	return warnings, errors
}

func validateAzureRMContainerRegistryWebhookServiceUri(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^https?://[^\s]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must start with http:// or https:// and must not contain whitespaces: %q", k, value))
	}

	return warnings, errors
}

func expandWebhookPropertiesCreateParameters(d *schema.ResourceData) *containerregistry.WebhookPropertiesCreateParameters {
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

func expandWebhookPropertiesUpdateParameters(d *schema.ResourceData) *containerregistry.WebhookPropertiesUpdateParameters {
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

func expandWebhookActions(d *schema.ResourceData) *[]containerregistry.WebhookAction {
	actions := make([]containerregistry.WebhookAction, 0)
	for _, action := range d.Get("actions").(*schema.Set).List() {
		actions = append(actions, containerregistry.WebhookAction(action.(string)))
	}

	return &actions
}
