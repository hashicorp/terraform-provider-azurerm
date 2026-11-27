// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_managed_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/managedcluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (k Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceId, err := managedcluster.ParseManagedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resourceID: %+v", err)
			}
			clusterClient := metadata.Client.ServiceFabricManaged.ManagedClusterClient

			if err = clusterClient.DeleteThenPoll(ctx, *resourceId); err != nil {
				return fmt.Errorf("while deleting cluster %q: %+v", resourceId.String(), err)
			}
			return nil
		},
		Timeout: 90 * time.Minute,
	}
}
