package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogicAppIntegrationAccountAssembly() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountAssemblyCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountAssemblyRead,
		Update: resourceLogicAppIntegrationAccountAssemblyCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountAssemblyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationAccountAssemblyID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountAssemblyArtifactName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"integration_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"assembly_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.IntegrationAccountAssemblyName(),
			},

			"assembly_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "0.0.0.0",
				ValidateFunc: validate.IntegrationAccountAssemblyVersion(),
			},

			"content": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"content", "content_link_uri"},
			},

			"content_link_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				AtLeastOneOf: []string{"content", "content_link_uri"},
			},

			"metadata": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceLogicAppIntegrationAccountAssemblyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountAssemblyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIntegrationAccountAssemblyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.AssemblyName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_assembly", id.ID())
		}
	}

	parameters := logic.AssemblyDefinition{
		Properties: &logic.AssemblyProperties{
			AssemblyName:    utils.String(d.Get("assembly_name").(string)),
			AssemblyVersion: utils.String(d.Get("assembly_version").(string)),
			ContentType:     utils.String("application/octet-stream"),
		},
	}

	if v, ok := d.GetOk("content"); ok {
		parameters.Properties.Content = utils.String(v.(string))
	}

	if v, ok := d.GetOk("content_link_uri"); ok {
		parameters.Properties.ContentLink = &logic.ContentLink{
			URI: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := v.(map[string]interface{})
		parameters.Properties.Metadata = &metadata
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IntegrationAccountName, id.AssemblyName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountAssemblyRead(d, meta)
}

func resourceLogicAppIntegrationAccountAssemblyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountAssemblyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountAssemblyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.AssemblyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AssemblyName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if props := resp.Properties; props != nil {
		d.Set("assembly_name", props.AssemblyName)
		d.Set("assembly_version", props.AssemblyVersion)
		d.Set("content_link_uri", d.Get("content_link_uri").(string))
		d.Set("content", d.Get("content").(string))

		if props.Metadata != nil {
			metadata := props.Metadata.(map[string]interface{})
			d.Set("metadata", metadata)
		}
	}

	return nil
}

func resourceLogicAppIntegrationAccountAssemblyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountAssemblyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountAssemblyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IntegrationAccountName, id.AssemblyName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
