// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity
type Resource struct{}

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_kubernetes_automatic_cluster"
}

func (r Resource) Identity() resourceids.ResourceId {
	return &commonids.KubernetesClusterId{}
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.Containers.KubernetesClustersClient
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if resp.Model.Sku == nil || resp.Model.Sku.Name == nil {
			return fmt.Errorf("importing %s: SKU information is missing", id)
		}

		if pointer.From(resp.Model.Sku.Name) != managedclusters.ManagedClusterSKUNameAutomatic {
			return fmt.Errorf("importing %s: specified Kubernetes Cluster is not using the SKU `Automatic`, got `%s`", id, pointer.From(resp.Model.Sku.Name))
		}

		return nil
	}
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateKubernetesClusterID
}
