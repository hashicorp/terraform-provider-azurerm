package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2018-07-12/botservice"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBotChannelSlack() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBotChannelSlackCreate,
		Read:   resourceArmBotChannelSlackRead,
		Delete: resourceArmBotChannelSlackDelete,
		Update: resourceArmBotChannelSlackUpdate,

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

			"verification_token": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"landing_page_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmBotChannelSlackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	botName := d.Get("bot_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, string(botservice.ChannelNameSlackChannel), botName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of creating Channel Slack for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_bot_channel_slack", *existing.ID)
		}
	}

	channel := botservice.BotChannel{
		Properties: botservice.SlackChannel{
			Properties: &botservice.SlackChannelProperties{
				ClientID:                utils.String(d.Get("client_id").(string)),
				ClientSecret:            utils.String(d.Get("client_secret").(string)),
				VerificationToken:       utils.String(d.Get("verification_token").(string)),
				LandingPageURL:          utils.String(d.Get("landing_page_url").(string)),
				IsEnabled:               utils.Bool(true),
				RegisterBeforeOAuthFlow: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameSlackChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Create(ctx, resourceGroup, botName, botservice.ChannelNameSlackChannel, channel); err != nil {
		return fmt.Errorf("Error issuing create request for Channel Slack for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		return fmt.Errorf("Error making get request for Channel Slack for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Channel Slack for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotChannelSlackRead(d, meta)
}

func resourceArmBotChannelSlackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	botName := id.Path["botServices"]
	resp, err := client.Get(ctx, id.ResourceGroup, botName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Channel Slack for Bot %q (Resource Group %q) was not found - removing from state!", id.ResourceGroup, botName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Channel Slack for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", resp.Location)
	d.Set("bot_name", botName)

	if props := resp.Properties; props != nil {
		if channel, ok := props.AsSlackChannel(); ok {
			if channelProps := channel.Properties; channelProps != nil {
				d.Set("client_id", channelProps.ClientID)
			}
		}
	}

	return nil
}

func resourceArmBotChannelSlackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	botName := d.Get("bot_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	channel := botservice.BotChannel{
		Properties: botservice.SlackChannel{
			Properties: &botservice.SlackChannelProperties{
				ClientID:                utils.String(d.Get("client_id").(string)),
				ClientSecret:            utils.String(d.Get("client_secret").(string)),
				VerificationToken:       utils.String(d.Get("verification_token").(string)),
				LandingPageURL:          utils.String(d.Get("landing_page_url").(string)),
				IsEnabled:               utils.Bool(true),
				RegisterBeforeOAuthFlow: utils.Bool(true),
			},
			ChannelName: botservice.ChannelNameSlackChannel1,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Kind:     botservice.KindBot,
	}

	if _, err := client.Update(ctx, resourceGroup, botName, botservice.ChannelNameSlackChannel, channel); err != nil {
		return fmt.Errorf("Error issuing create request for Channel Slack for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, botName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		return fmt.Errorf("Error making get request for Channel Slack for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Channel Slack for Bot %q (Resource Group %q): %+v", resourceGroup, botName, err)
	}

	d.SetId(*resp.ID)

	return resourceArmBotChannelSlackRead(d, meta)
}

func resourceArmBotChannelSlackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Bot.ChannelClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	botName := id.Path["botServices"]

	resp, err := client.Delete(ctx, id.ResourceGroup, botName, string(botservice.ChannelNameSlackChannel))
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Channel Slack for Bot %q (Resource Group %q): %+v", id.ResourceGroup, botName, err)
		}
	}

	return nil
}
