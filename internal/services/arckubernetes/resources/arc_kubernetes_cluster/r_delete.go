// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	arckubernetesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArcKubernetesClusterDelete(d *pluginsdk.ResourceData, meta any) error {
	client := arckubernetesClient.From(meta).ArcKubernetesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := connectedclusters.ParseConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	if err := client.ConnectedClusterDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
