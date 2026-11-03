// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_service

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareServiceCreateUpdate,
		Read:   resourceHealthcareServiceRead,
		Update: resourceHealthcareServiceCreateUpdate,
		Delete: resourceHealthcareServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := resource.ParseServiceID(id)
			return err
		}),

		Schema: healthcareServiceSchema(),
	}
}
