package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/bot/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBotConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBotConnectionCreate,
		Read:   resourceArmBotConnectionRead,
		Update: resourceArmBotConnectionUpdate,
		Delete: resourceArmBotConnectionDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.BotConnectionID(id)
			return err
		}),

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"bot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"service_provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_secret": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"scopes": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmBotConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewBotConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("bot_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ConnectionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Bot Connection %q (Bot %q / Resource Group %q): %+v", resourceId.ConnectionName, resourceId.BotServiceName, resourceId.ResourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_connection", resourceId.ID())
		}
	}

	serviceProviderName := d.Get("service_provider_name").(string)
	var serviceProviderId *string
	var availableProviders []string

	serviceProviders, err := client.ListServiceProviders(ctx)
	if err != nil {
		return fmt.Errorf("listing Bot Connection service provider: %+v", err)
	}

	if serviceProviders.Value == nil {
		return fmt.Errorf("no service providers were returned from the Azure API")
	}
	for _, provider := range *serviceProviders.Value {
		if provider.Properties == nil || provider.Properties.ServiceProviderName == nil {
			continue
		}
		name := provider.Properties.ServiceProviderName
		if strings.EqualFold(serviceProviderName, *name) {
			serviceProviderId = provider.Properties.ID
			break
		}
		availableProviders = append(availableProviders, *name)
	}

	if serviceProviderId == nil {
		return fmt.Errorf("the Service Provider %q was not found. The available service providers are %s", serviceProviderName, strings.Join(availableProviders, ","))
	}

	connection := botservice.ConnectionSetting{
		Properties: &botservice.ConnectionSettingProperties{
			ServiceProviderID: serviceProviderId,
			ClientID:          utils.String(d.Get("client_id").(string)),
			ClientSecret:      utils.String(d.Get("client_secret").(string)),
			Scopes:            utils.String(d.Get("scopes").(string)),
			Parameters:        expandBotConnectionParameters(d.Get("parameters").(map[string]interface{})),
		},
		Kind:     botservice.KindBot,
		Location: utils.String(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.BotServiceName, resourceId.ConnectionName, connection); err != nil {
		return fmt.Errorf("creating Bot Connection %q (Bot %q / Resource Group %q): %+v", resourceId.ConnectionName, resourceId.BotServiceName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceArmBotConnectionRead(d, meta)
}

func resourceArmBotConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Bot Connection %q (Bot %q / Resource Group %q) was not found - removing from state!", id.ConnectionName, id.BotServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Bot Connection %q (Bot %q / Resource Group %q): %+v", id.ConnectionName, id.BotServiceName, id.ResourceGroup, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("bot_name", id.BotServiceName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.Properties; props != nil {
		d.Set("client_id", props.ClientID)
		d.Set("scopes", props.Scopes)
		if err := d.Set("parameters", flattenBotConnectionParameters(props.Parameters)); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmBotConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotConnectionID(d.Id())
	if err != nil {
		return err
	}

	connection := botservice.ConnectionSetting{
		Properties: &botservice.ConnectionSettingProperties{
			ServiceProviderDisplayName: utils.String(d.Get("service_provider_name").(string)),
			ClientID:                   utils.String(d.Get("client_id").(string)),
			ClientSecret:               utils.String(d.Get("client_secret").(string)),
			Scopes:                     utils.String(d.Get("scopes").(string)),
			Parameters:                 expandBotConnectionParameters(d.Get("parameters").(map[string]interface{})),
		},
		Kind:     botservice.KindBot,
		Location: utils.String(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName, connection); err != nil {
		return fmt.Errorf("updating Bot Connection %q (Bot %q / Resource Group %q): %+v", id.ConnectionName, id.BotServiceName, id.ResourceGroup, err)
	}

	return resourceArmBotConnectionRead(d, meta)
}

func resourceArmBotConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BotConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.BotServiceName, id.ConnectionName)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Bot Connection %q (Bot %q / Resource Group %q): %+v", id.ConnectionName, id.BotServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandBotConnectionParameters(input map[string]interface{}) *[]botservice.ConnectionSettingParameter {
	output := make([]botservice.ConnectionSettingParameter, 0)

	for k, v := range input {
		output = append(output, botservice.ConnectionSettingParameter{
			Key:   utils.String(k),
			Value: utils.String(v.(string)),
		})
	}
	return &output
}

func flattenBotConnectionParameters(input *[]botservice.ConnectionSettingParameter) map[string]interface{} {
	output := make(map[string]interface{})
	if input == nil {
		return output
	}

	for _, parameter := range *input {
		if key := parameter.Key; key != nil {
			// We disregard the clientSecret and clientId as one is sensitive and the other is returned in the ClientId attribute.
			if *key != "clientSecret" && *key != "clientId" && *key != "scopes" {
				if value := parameter.Value; value != nil {
					output[*key] = *value
				}
			}
		}
	}

	return output
}
