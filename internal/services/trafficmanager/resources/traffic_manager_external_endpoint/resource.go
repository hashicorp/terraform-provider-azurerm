// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_external_endpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExternalEndpointCreate,
		Read:   resourceExternalEndpointRead,
		Update: resourceExternalEndpointUpdate,
		Delete: resourceExternalEndpointDelete,
		Importer: schema.ValidateResourceIDPriorToImport(func(id string) error {
			endpointType, err := trafficmanagers.ParseEndpointTypeID(id)
			if err != nil {
				return err
			}

			if endpointType.EndpointType != trafficmanagers.EndpointTypeExternalEndpoints {
				return fmt.Errorf("this resource only supports `ExternalEndpoints` but got %s", string(endpointType.EndpointType))
			}

			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: trafficManagerExternalEndpointSchema(),
	}
}
