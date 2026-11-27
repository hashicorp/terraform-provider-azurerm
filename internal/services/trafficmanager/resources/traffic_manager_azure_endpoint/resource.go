// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_azure_endpoint

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAzureEndpointCreate,
		Read:   resourceAzureEndpointRead,
		Update: resourceAzureEndpointUpdate,
		Delete: resourceAzureEndpointDelete,
		Importer: schema.ValidateResourceIDPriorToImport(func(id string) error {
			endpointType, err := trafficmanagers.ParseEndpointTypeID(id)
			if err != nil {
				return err
			}

			if endpointType.EndpointType != trafficmanagers.EndpointTypeAzureEndpoints {
				return fmt.Errorf("this resource only supports `AzureEndpoints` but got %s", string(endpointType.EndpointType))
			}

			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: trafficManagerAzureEndpointSchema(),
	}
}
