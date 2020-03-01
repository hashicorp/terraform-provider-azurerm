package eventgrid

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmEventGridDomainTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventGridDomainTopicCreateUpdate,
		Read:   resourceArmEventGridDomainTopicRead,
		Update: resourceArmEventGridDomainTopicCreateUpdate,
		Delete: resourceArmEventGridDomainTopicDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),
		},
	}
}

func resourceArmEventGridDomainTopicCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.DomainTopicsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	domainName := d.Get("domain_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, domainName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventGrid Domain Topic %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventgrid_domain_topic", *existing.ID)
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, domainName, name)
	if err != nil {
		return fmt.Errorf("Error creating/updating EventGrid Domain Topic %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for EventGrid Domain Topic %q (Resource Group %q) to become available: %s", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, domainName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving EventGrid Domain Topic %q (Resource Group %q): %s", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read EventGrid Domain Topic %q (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventGridDomainTopicRead(d, meta)
}

func resourceArmEventGridDomainTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.DomainTopicsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	domainName := id.Path["domains"]
	name := id.Path["topics"]

	resp, err := client.Get(ctx, resourceGroup, domainName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] EventGrid Domain Topic %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on EventGrid Domain Topic %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("domain_name", domainName)
	d.Set("resource_group_name", resourceGroup)

	return nil
}

func resourceArmEventGridDomainTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.DomainTopicsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	domainName := id.Path["domains"]
	name := id.Path["topics"]

	future, err := client.Delete(ctx, resGroup, domainName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Domain Topic %q: %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Domain Topic %q: %+v", name, err)
	}

	return nil
}
