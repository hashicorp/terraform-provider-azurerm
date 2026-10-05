// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package image

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceImageCreateUpdate,
		Read:   resourceImageRead,
		Update: resourceImageCreateUpdate,
		Delete: resourceImageDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := images.ParseImageID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: imageSchema(),
	}
}
