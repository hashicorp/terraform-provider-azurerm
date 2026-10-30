// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_certificate

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/dpscertificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubDPSCertificateCreate,
		Read:   resourceIotHubDPSCertificateRead,
		Update: resourceIotHubDPSCertificateUpdate,
		Delete: resourceIotHubDPSCertificateDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dpscertificate.ParseCertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: iothubDpsCertificateSchema(),
	}
}
