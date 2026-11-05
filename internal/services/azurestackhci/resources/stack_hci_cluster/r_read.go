// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_cluster

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmStackHCIClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	hciAssignmentsClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient
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
		d.Set("identity", flattenSystemAssigned(model.Identity))

		if props := model.Properties; props != nil {
			d.Set("client_id", props.AadClientId)
			d.Set("tenant_id", props.AadTenantId)
			d.Set("cloud_id", props.CloudId)
			d.Set("service_endpoint", props.ServiceEndpoint)
			d.Set("resource_provider_object_id", props.ResourceProviderObjectId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	hclAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")
	assignmentResp, err := hciAssignmentsClient.Get(ctx, hclAssignmentId)
	if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
		return err
	}
	configId := ""
	if model := assignmentResp.Model; model != nil && model.Properties != nil && model.Properties.ConfigurationProfile != nil {
		parsed, err := configurationprofiles.ParseConfigurationProfileIDInsensitively(*model.Properties.ConfigurationProfile)
		if err != nil {
			return err
		}
		configId = parsed.ID()
	}
	d.Set("automanage_configuration_id", configId)

	return nil
}

func flattenSystemAssigned(input *identity.SystemAndUserAssignedMap) []any {
	if input == nil {
		return []any{}
	}

	if input.Type == identity.TypeNone {
		return []any{}
	}

	return []any{
		map[string]any{
			"type":         input.Type,
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}
}
