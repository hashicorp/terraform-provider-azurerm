// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	arckubernetesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArcKubernetesClusterUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := arckubernetesClient.From(meta).ArcKubernetesClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := connectedclusters.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	props := connectedclusters.ConnectedClusterPatch{
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, err := client.ConnectedClusterUpdate(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArcKubernetesClusterRead(d, meta)
}
