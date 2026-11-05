// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArcKubernetesClusterUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ArcKubernetes.ArcKubernetesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
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
