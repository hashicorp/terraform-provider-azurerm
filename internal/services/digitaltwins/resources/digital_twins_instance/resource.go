// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_instance

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDigitalTwinsInstanceCreate,
		Read:   resourceDigitalTwinsInstanceRead,
		Update: resourceDigitalTwinsInstanceUpdate,
		Delete: resourceDigitalTwinsInstanceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(id)
			return err
		}),

		Schema: digitalTwinsInstanceSchema(),
	}
}
