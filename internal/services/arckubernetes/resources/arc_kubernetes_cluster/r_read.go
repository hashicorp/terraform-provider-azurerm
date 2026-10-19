// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArcKubernetesClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connectedclusters.ParseConnectedClusterID(d.Id())
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
