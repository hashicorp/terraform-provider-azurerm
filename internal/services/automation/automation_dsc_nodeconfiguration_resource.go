package automation

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			_, err := parse.NodeConfigurationID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"content_embedded": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"content_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"content_hash": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"value": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"configuration_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"increment_build_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceAutomationDscNodeConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscNodeConfigurationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Dsc Node Configuration creation.")

	id := parse.NewNodeConfigurationID(client.SubscriptionID, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automation_dsc_nodeconfiguration", id.ID())
		}
	}

	content := d.Get("content_embedded").(string)
	contentVersion := d.Get("content_version").(string)

	// configuration name is always the first part of the dsc node configuration
	// e.g. webserver.prod or webserver.local will be associated to the dsc configuration webserver

	configurationName := strings.Split(id.Name, ".")[0]

	tagsVal := tags.Expand(d.Get("tags").(map[string]interface{}))

	parameters := automation.DscNodeConfigurationCreateOrUpdateParameters{
		DscNodeConfigurationCreateOrUpdateParametersProperties: &automation.DscNodeConfigurationCreateOrUpdateParametersProperties{
			Source: &automation.ContentSource{
				Type:    automation.ContentSourceTypeEmbeddedContent,
				Value:   utils.String(content),
				Version: utils.String(contentVersion),
				Hash:    expandContentHash(d.Get("content_hash").([]interface{})),
			},
			Configuration: &automation.DscConfigurationAssociationProperty{
				Name: utils.String(configurationName),
			},
			IncrementNodeConfigurationBuild: utils.Bool(d.Get("increment_build_enabled").(bool)),
		},
		Tags: tagsVal,
		Name: utils.String(id.Name),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update for %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAutomationDscNodeConfigurationRead(d, meta)
}

func resourceAutomationDscNodeConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscNodeConfigurationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NodeConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on AzureRM Automation Dsc Node Configuration %q: %+v", id.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("automation_account_name", id.AutomationAccountName)
	d.Set("configuration_name", resp.Configuration.Name)

	// cannot read back content_embedded, tags, content_version, increment_build_enabled, hash_xx as not part of body nor exposed through method

	return nil
}

func resourceAutomationDscNodeConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.DscNodeConfigurationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NodeConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("issuing AzureRM delete request for Automation Dsc Node Configuration %q: %+v", id.Name, err)
	}

	return nil
}

func expandContentHash(input []interface{}) *automation.ContentHash {
	if len(input) == 0 {
		return nil
	}
	meta := input[0].(map[string]interface{})
	var hash automation.ContentHash
	hash.Algorithm = utils.String(meta["algorithm"].(string))
	hash.Value = utils.String(meta["value"].(string))
	return &hash
}
