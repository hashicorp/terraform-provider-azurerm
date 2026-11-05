// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster_extension

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2024-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_arc_kubernetes_cluster_extension"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(val any, key string) (warns []string, errs []error) {
		idRaw, ok := val.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("expected `id` to be a string but got %+v", val))
			return
		}

		id, err := extensions.ParseScopedExtensionID(idRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("parsing %q: %+v", idRaw, err))
			return
		}

		// validate the scope is a connected cluster id
		if _, err := connectedclusters.ParseConnectedClusterID(id.Scope); err != nil {
			errs = append(errs, fmt.Errorf("parsing %q as a Connected Cluster ID: %+v", idRaw, err))
			return
		}

		return
	}
}
