// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.Containers.KubernetesClustersClient

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := clusterClient.DeleteThenPoll(ctx, *id, managedclusters.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %w", *id, err)
			}

			return nil
		},
	}
}
