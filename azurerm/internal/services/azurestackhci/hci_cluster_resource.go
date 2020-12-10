package azurestackhci

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/azurestackhci/mgmt/2020-10-01/azurestackhci"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/azurestackhci/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/azurestackhci/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHCICluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHCIClusterCreate,
		Read:   resourceArmHCIClusterRead,
		Update: resourceArmHCIClusterUpdate,
		Delete: resourceArmHCIClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ClusterID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"client_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceArmHCIClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.ClusterClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewClusterID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing HCI Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_hci_cluster", *existing.ID)
	}

	cluster := azurestackhci.Cluster{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		ClusterProperties: &azurestackhci.ClusterProperties{
			AadClientID: utils.String(d.Get("client_id").(string)),
			AadTenantID: utils.String(d.Get("tenant_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Create(ctx, resourceGroup, name, cluster); err != nil {
		return fmt.Errorf("creating HCI Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID(""))

	return resourceArmHCIClusterRead(d, meta)
}

func resourceArmHCIClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] HCI Cluster %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HCI Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.ClusterProperties; props != nil {
		d.Set("client_id", props.AadClientID)
		d.Set("tenant_id", props.AadTenantID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmHCIClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	cluster := azurestackhci.ClusterUpdate{}

	if d.HasChange("tags") {
		cluster.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, cluster); err != nil {
		return fmt.Errorf("updating HCI Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceArmHCIClusterRead(d, meta)
}

func resourceArmHCIClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting HCI Cluster %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
