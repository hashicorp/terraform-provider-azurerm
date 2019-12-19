package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBotChannelEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBotChannelEmailCreate,
		Read:   resourceArmBotChannelEmailRead,
		Delete: resourceArmBotChannelEmailDelete,
		Update: resourceArmBotChannelEmailUpdate,

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
			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"bot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"email_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"email_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmBotChannelEmailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	botName := d.Get("bot_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, string(botservice.ChannelNameEmailChannel), botName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Channel Email for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_channel_email", *existing.ID)
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.EmailChannel{
			Properties: &botservice.EmailChannelProperties{
				EmailAddress: utils.String(d.Get("email_address").(string)),
				Password:     utils.String(d.Get("email_password").(string)),
				IsEnabled:    utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameEmailChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, resourceGroup, botName, botservice.ChannelNameEmailChannel, channel); err != nil {
		return fmt.Errorf("Error issuing create request for Channel Email for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameEmailChannel))
	if err != nil {
		return fmt.Errorf("Error making get request for Channel Email for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Channel Email for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotChannelEmailRead(d, meta)
}

func resourceArmBotChannelEmailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	botName := id.Path["botServices"]
	resp, err := client.Get(ctx, id.ResourceGroup, botName, string(botservice.ChannelNameEmailChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Channel Email for Bot %q (Resource Group %q) was not found - removing from state!", id.ResourceGroup, botName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Channel Email for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", resp.Location)
	d.Set("bot_name", botName)

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsEmailChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("email_address", channelProps.EmailAddress)
			}
		}
	}

	return nil
}

func resourceArmBotChannelEmailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	botName := id.Path["botServices"]

	channel := botservice.BotChannel{
		Properties: botservice.EmailChannel{
			Properties: &botservice.EmailChannelProperties{
				EmailAddress: utils.String(d.Get("email_address").(string)),
				Password:     utils.String(d.Get("email_password").(string)),
				IsEnabled:    utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameEmailChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, id.ResourceGroup, botName, botservice.ChannelNameEmailChannel, channel); err != nil {
		return fmt.Errorf("Error issuing create request for Channel Email for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, botName, string(botservice.ChannelNameEmailChannel))
	if err != nil {
		return fmt.Errorf("Error making get request for Channel Email for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Channel Email for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotChannelEmailRead(d, meta)
}

func resourceArmBotChannelEmailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	botName := id.Path["botServices"]

	resp, err := client.Delete(ctx, id.ResourceGroup, botName, string(botservice.ChannelNameEmailChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Channel Email for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
		}
	}

	return nil
}
