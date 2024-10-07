// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/dscconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationDscConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationDscConfigurationCreateUpdate,
		Read:   resourceAutomationDscConfigurationRead,
		Update: resourceAutomationDscConfigurationCreateUpdate,
		Delete: resourceAutomationDscConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dscconfiguration.ParseConfigurationID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z0-9_]{1,64}$`),
					`The name length must be from 1 to 64 characters. The name can only contain letters, numbers and underscores.`,
				),
			},

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"content_embedded": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"log_verbose": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceAutomationDscConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscConfiguration
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Dsc Configuration creation.")

	id := dscconfiguration.NewConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_dsc_configuration", id.ID())
		}
	}

	contentEmbedded := d.Get("content_embedded").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	logVerbose := d.Get("log_verbose").(bool)
	description := d.Get("description").(string)

	parameters := dscconfiguration.DscConfigurationCreateOrUpdateParameters{
		Properties: dscconfiguration.DscConfigurationCreateOrUpdateProperties{
			LogVerbose:  utils.Bool(logVerbose),
			Description: utils.String(description),
			Source: dscconfiguration.ContentSource{
				Type:  pointer.To(dscconfiguration.ContentSourceTypeEmbeddedContent),
				Value: utils.String(contentEmbedded),
			},
		},
		Location: utils.String(location),
		Tags:     pointer.To(expandStringInterfaceMap(d.Get("tags").(map[string]interface{}))),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceAutomationDscConfigurationRead(d, meta)
}

func resourceAutomationDscConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscConfiguration
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dscconfiguration.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.ConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if props := model.Properties; props != nil {
			d.Set("log_verbose", props.LogVerbose)
			d.Set("description", props.Description)
			d.Set("state", string(pointer.From(props.State)))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	contentResp, err := client.GetContent(ctx, *id)
	if err != nil {
		return fmt.Errorf("making Read request on AzureRM Automation Dsc Configuration content %q: %+v", id.ConfigurationName, err)
	}

	if contentHttpResponse := contentResp.HttpResponse; contentHttpResponse != nil {
		if contentHttpResponse.Body != nil {
			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(contentResp.HttpResponse.Body); err != nil {
				return fmt.Errorf("reading from AzureRM Automation Dsc Configuration buffer %q: %+v", id.ConfigurationName, err)
			}
			content := buf.String()

			d.Set("content_embedded", content)
		}
	}

	return nil
}

func resourceAutomationDscConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscConfiguration
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dscconfiguration.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
