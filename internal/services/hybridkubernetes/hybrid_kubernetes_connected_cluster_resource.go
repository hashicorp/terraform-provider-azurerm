package hybridkubernetes

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hybridkubernetes/sdk/2021-10-01/hybridkubernetes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHybridKubernetesConnectedCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHybridKubernetesConnectedClusterCreate,
		Read:   resourceHybridKubernetesConnectedClusterRead,
		Update: resourceHybridKubernetesConnectedClusterUpdate,
		Delete: resourceHybridKubernetesConnectedClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := hybridkubernetes.ParseConnectedClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"agent_public_key_certificate": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"distribution": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"identity": commonschema.SystemAssignedIdentityRequired(),

			"infrastructure": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceHybridKubernetesConnectedClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HybridKubernetes.HybridKubernetesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := hybridkubernetes.NewConnectedClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.ConnectedClusterGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_hybrid_kubernetes_connected_cluster", id.ID())
		}
	}

	identity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	location := azure.NormalizeLocation(d.Get("location"))
	props := hybridkubernetes.ConnectedCluster{
		Identity: *identity,
		Location: location,
		Properties: hybridkubernetes.ConnectedClusterProperties{
			AgentPublicKeyCertificate: d.Get("agent_public_key_certificate").(string),
			Distribution:              utils.String(d.Get("distribution").(string)),
			Infrastructure:            utils.String(d.Get("infrastructure").(string)),
		},
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.ConnectedClusterCreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceHybridKubernetesConnectedClusterRead(d, meta)
}

func resourceHybridKubernetesConnectedClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HybridKubernetes.HybridKubernetesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hybridkubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ConnectedClusterGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		if err := d.Set("identity", identity.FlattenSystemAssigned(&model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("location", location.Normalize(model.Location))
		props := model.Properties
		d.Set("agent_public_key_certificate", props.AgentPublicKeyCertificate)
		d.Set("distribution", props.Distribution)
		d.Set("infrastructure", props.Infrastructure)
		d.Set("provisioning_state", props.ProvisioningState)

		if err := tagsHelper.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceHybridKubernetesConnectedClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HybridKubernetes.HybridKubernetesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hybridkubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	props := hybridkubernetes.ConnectedClusterPatch{
		Tags: tagsHelper.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.ConnectedClusterUpdate(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceHybridKubernetesConnectedClusterRead(d, meta)
}

func resourceHybridKubernetesConnectedClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HybridKubernetes.HybridKubernetesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hybridkubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	if err := client.ConnectedClusterDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
