package azurestackhci

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2020-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmStackHCICluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmStackHCIClusterCreate,
		Read:   resourceArmStackHCIClusterRead,
		Update: resourceArmStackHCIClusterUpdate,
		Delete: resourceArmStackHCIClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"client_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceArmStackHCIClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.ClusterClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := clusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_stack_hci_cluster", id.ID())
	}

	cluster := clusters.Cluster{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &clusters.ClusterProperties{
			AadClientId: d.Get("client_id").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		cluster.Properties.AadTenantId = v.(string)
	} else {
		tenantId := meta.(*clients.Client).Account.TenantId
		cluster.Properties.AadTenantId = tenantId
	}

	if _, err := client.Create(ctx, id, cluster); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmStackHCIClusterRead(d, meta)
}

func resourceArmStackHCIClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.ClusterClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("client_id", props.AadClientId)
			d.Set("tenant_id", props.AadTenantId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceArmStackHCIClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.ClusterClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	cluster := clusters.ClusterUpdate{}

	if d.HasChange("tags") {
		cluster.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, cluster); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArmStackHCIClusterRead(d, meta)
}

func resourceArmStackHCIClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.ClusterClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
