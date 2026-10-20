// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_flux_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2025-04-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_kubernetes_flux_configuration"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(val any, key string) (warns []string, errs []error) {
		idRaw, ok := val.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("expected `id` to be a string but got %+v", val))
			return
		}

		id, err := fluxconfiguration.ParseScopedFluxConfigurationID(idRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("parsing %q: %+v", idRaw, err))
			return
		}

		// validate the scope is a connected cluster id
		if _, err := commonids.ParseKubernetesClusterID(id.Scope); err != nil {
			errs = append(errs, fmt.Errorf("parsing %q as a Kubernetes Cluster ID: %+v", idRaw, err))
			return
		}

		return
	}
}
