// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmStackHCIClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := clusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_stack_hci_cluster", id.ID())
		}
	}

	cluster := clusters.Cluster{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &clusters.ClusterProperties{
			AadClientId: pointer.To(d.Get("client_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("identity"); ok {
		cluster.Identity = expandSystemAssigned(v.([]any))
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		cluster.Properties.AadTenantId = pointer.To(v.(string))
	} else {
		tenantId := meta.(*clients.Client).Account.TenantId
		cluster.Properties.AadTenantId = pointer.To(tenantId)
	}

	if _, err := client.Create(ctx, id, cluster); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if v, ok := d.GetOk("automanage_configuration_id"); ok {
		configurationProfilesClient := meta.(*clients.Client).Automanage.ConfigurationProfilesClient
		hciAssignmentsClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient

		configurationProfileId, err := configurationprofiles.ParseConfigurationProfileID(v.(string))
		if err != nil {
			return err
		}

		if _, err = configurationProfilesClient.Get(ctx, *configurationProfileId); err != nil {
			return fmt.Errorf("checking for existing %s: %+v", configurationProfileId, err)
		}

		hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(subscriptionId, id.ResourceGroupName, id.ClusterName, "default")
		assignmentsResp, err := hciAssignmentsClient.Get(ctx, hciAssignmentId)
		if err != nil && !response.WasNotFound(assignmentsResp.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", hciAssignmentId, err)
		}

		if response.WasNotFound(assignmentsResp.HttpResponse) {
			properties := configurationprofilehciassignments.ConfigurationProfileAssignment{
				Properties: &configurationprofilehciassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: pointer.To(configurationProfileId.ID()),
				},
			}

			if _, err := hciAssignmentsClient.CreateOrUpdate(ctx, hciAssignmentId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", hciAssignmentId, err)
			}
		}
	}

	return resourceArmStackHCIClusterRead(d, meta)
}

// API does not accept userAssignedIdentity as in swagger https://github.com/Azure/azure-rest-api-specs/issues/28260
func expandSystemAssigned(input []any) *identity.SystemAndUserAssignedMap {
	if len(input) == 0 || input[0] == nil {
		return &identity.SystemAndUserAssignedMap{
			Type: identity.TypeNone,
		}
	}

	return &identity.SystemAndUserAssignedMap{
		Type: identity.TypeSystemAssigned,
	}
}
