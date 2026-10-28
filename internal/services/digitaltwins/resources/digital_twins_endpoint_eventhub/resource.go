// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_endpoint_eventhub

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDigitalTwinsEndpointEventHubCreateUpdate,
		Read:   resourceDigitalTwinsEndpointEventHubRead,
		Update: resourceDigitalTwinsEndpointEventHubCreateUpdate,
		Delete: resourceDigitalTwinsEndpointEventHubDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := endpoints.ParseEndpointID(id)
			return err
		}, helpers.ValidateEndpointType(func(input endpoints.DigitalTwinsEndpointResourceProperties) error {
			if _, ok := input.(endpoints.EventHub); !ok {
				return fmt.Errorf("expected an EventHub type but got: %+v", input)
			}
			return nil
		})),

		Schema: digitalTwinsEndpointEventhubSchema(),
	}
}
