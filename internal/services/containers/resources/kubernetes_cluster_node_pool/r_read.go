// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_node_pool

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKubernetesClusterNodePoolRead(d *pluginsdk.ResourceData, meta any) error {
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := poolsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceKubernetesClusterNodePoolFlatten(d, id, resp.Model)
}

func upgradeSettingsSchemaNodePoolResource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_surge": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					ValidateFunc:  validation.StringIsNotEmpty,
					ConflictsWith: []string{"upgrade_settings.0.max_unavailable"},
				},
				"drain_timeout_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},
				"max_unavailable": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					ValidateFunc:  validation.StringIsNotEmpty,
					ConflictsWith: []string{"upgrade_settings.0.max_surge"},
				},
				"node_soak_duration_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 30),
				},
				"undrainable_node_behavior": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForUndrainableNodeBehavior(), true),
				},
			},
		},
	}
}
