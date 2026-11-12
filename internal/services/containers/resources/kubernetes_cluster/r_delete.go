// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKubernetesClusterDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKubernetesClusterID(d.Id())
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		if _, err := client.Delete(ctx, maintenanceId); err != nil {
			return fmt.Errorf("deleting Maintenance Configuration for %s: %+v", *id, err)
		}
	}

	if err = client.DeleteThenPoll(ctx, *id, managedclusters.DefaultDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
