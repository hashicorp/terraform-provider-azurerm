// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arckubernetes

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	arckubernetes "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArcKubernetesCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArcKubernetesClusterCreate,
		Read:   resourceArcKubernetesClusterRead,
		Update: resourceArcKubernetesClusterUpdate,
		Delete: resourceArcKubernetesClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := arckubernetes.ParseConnectedClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-_a-zA-Z0-9]{1,260}$"),
					"The name of Arc Kubernetes Cluster can only include alphanumeric characters, underscores, hyphens, has a maximum length of 260 characters, and must be unique.",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"agent_public_key_certificate": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.Base64EncodedString,
			},

			"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

			"location": commonschema.Location(),

			"agent_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"distribution": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"infrastructure": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kubernetes_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"offering": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"total_core_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"total_node_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceArcKubernetesClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := arckubernetes.NewConnectedClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.ConnectedClusterGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_arc_kubernetes_cluster", id.ID())
	}

	identityValue, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	location := location.Normalize(d.Get("location").(string))
	props := arckubernetes.ConnectedCluster{
		Identity: *identityValue,
		Location: location,
		Properties: arckubernetes.ConnectedClusterProperties{
			AgentPublicKeyCertificate: d.Get("agent_public_key_certificate").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.ConnectedClusterCreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArcKubernetesClusterRead(d, meta)
}

func resourceArcKubernetesClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := arckubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ConnectedClusterGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectedClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		if err := d.Set("identity", identity.FlattenSystemAssigned(&model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("location", location.Normalize(model.Location))
		props := model.Properties
		d.Set("agent_public_key_certificate", props.AgentPublicKeyCertificate)
		d.Set("agent_version", props.AgentVersion)
		d.Set("distribution", props.Distribution)
		d.Set("infrastructure", props.Infrastructure)
		d.Set("kubernetes_version", props.KubernetesVersion)
		d.Set("offering", props.Offering)
		d.Set("total_core_count", props.TotalCoreCount)
		d.Set("total_node_count", props.TotalNodeCount)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceArcKubernetesClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := arckubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	props := arckubernetes.ConnectedClusterPatch{
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.ConnectedClusterUpdate(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArcKubernetesClusterRead(d, meta)
}

func resourceArcKubernetesClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := arckubernetes.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	if err := client.ConnectedClusterDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
