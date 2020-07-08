package eventhub

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmEventHubCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventHubClusterCreateUpdate,
		Read:   resourceArmEventHubClusterRead,
		Update: resourceArmEventHubClusterCreateUpdate,
		Delete: resourceArmEventHubClusterDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ClusterID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			// You can't delete a cluster until at least 4 hours have passed from the initial creation.
			Delete: schema.DefaultTimeout(300 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Dedicated_1",
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmEventHubClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM EventHub Cluster creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	cluster := eventhub.Cluster{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku:      expandEventHubClusterSkuName(d.Get("sku_name").(string)),
	}

	future, err := client.Put(ctx, resourceGroup, name, cluster)
	if err != nil {
		return fmt.Errorf("creating EventHub Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of EventHub Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("making Read request on Azure EventHub Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("cannot read EventHub Cluster %s (Resource Group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventHubClusterRead(d, meta)
}

func resourceArmEventHubClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure EventHub Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("sku_name", flattenEventHubClusterSkuName(resp.Sku))
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmEventHubClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	// The EventHub Cluster can't be deleted until four hours after creation so we'll keep retrying until it can be deleted.
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if response.WasNotFound(future.Response()) {
				return nil
			}
			if strings.Contains(err.Error(), "Cluster cannot be deleted until four hours after its creation time") || future.Response().StatusCode == 429 {
				return resource.RetryableError(fmt.Errorf("expected eventhub cluster to be deleted but was in pending creation state, retrying"))
			}
			return resource.NonRetryableError(fmt.Errorf("issuing delete request for EventHub Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err))
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return resource.NonRetryableError(fmt.Errorf("deleting EventHub Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err))
		}

		return nil
	})
}

func expandEventHubClusterSkuName(skuName string) *eventhub.ClusterSku {
	if len(skuName) == 0 {
		return nil
	}

	name, capacity, err := azure.SplitSku(skuName)
	if err != nil {
		return nil
	}

	return &eventhub.ClusterSku{
		Name:     utils.String(name),
		Capacity: utils.Int32(capacity),
	}
}

func flattenEventHubClusterSkuName(input *eventhub.ClusterSku) string {
	if input == nil || input.Name == nil {
		return ""
	}

	return fmt.Sprintf("%s_%d", *input.Name, *input.Capacity)
}
