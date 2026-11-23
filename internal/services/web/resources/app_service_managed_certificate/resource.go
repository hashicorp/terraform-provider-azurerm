// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_managed_certificate

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceManagedCertificateCreate,
		Read:   resourceAppServiceManagedCertificateRead,
		Update: resourceAppServiceManagedCertificateUpdate,
		Delete: resourceAppServiceManagedCertificateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := certificates.ParseCertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: appServiceManagedCertificateSchema(),
	}
}
