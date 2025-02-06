// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/dscnodeconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationDscNodeConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationDscNodeConfigurationCreateUpdate,
		Read:   resourceAutomationDscNodeConfigurationRead,
		Update: resourceAutomationDscNodeConfigurationCreateUpdate,
		Delete: resourceAutomationDscNodeConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dscnodeconfiguration.ParseNodeConfigurationID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"content_embedded": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"configuration_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAutomationDscNodeConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscNodeConfiguration
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Dsc Node Configuration creation.")

	id := dscnodeconfiguration.NewNodeConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_dsc_nodeconfiguration", id.ID())
		}
	}

	content := d.Get("content_embedded").(string)

	// configuration name is always the first part of the dsc node configuration
	// e.g. webserver.prod or webserver.local will be associated to the dsc configuration webserver

	configurationName := strings.Split(id.NodeConfigurationName, ".")[0]

	contentSourceType := dscnodeconfiguration.ContentSourceTypeEmbeddedContent

	parameters := dscnodeconfiguration.DscNodeConfigurationCreateOrUpdateParameters{
		Properties: &dscnodeconfiguration.DscNodeConfigurationCreateOrUpdateParametersProperties{
			Source: dscnodeconfiguration.ContentSource{
				Type:  &contentSourceType,
				Value: utils.String(content),
			},
			Configuration: dscnodeconfiguration.DscConfigurationAssociationProperty{
				Name: utils.String(configurationName),
			},
		},
		Name: utils.String(id.NodeConfigurationName),
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationDscNodeConfigurationRead(d, meta)
}

func resourceAutomationDscNodeConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscNodeConfiguration
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dscnodeconfiguration.ParseNodeConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NodeConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			configurationName := ""
			if props.Configuration != nil && props.Configuration.Name != nil {
				configurationName = *props.Configuration.Name
			}
			d.Set("configuration_name", configurationName)
		}
	}

	// cannot read back content_embedded as not part of body nor exposed through method

	return nil
}

func resourceAutomationDscNodeConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscNodeConfiguration
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dscnodeconfiguration.ParseNodeConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
