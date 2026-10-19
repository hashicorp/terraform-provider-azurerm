// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_shared_access_policy

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// TODO: outside of this pr make this private

var IothubResourceName = "azurerm_iothub"

type accessRights struct {
	registryRead   bool
	registryWrite  bool
	serviceConnect bool
	deviceConnect  bool
}

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubDPSSharedAccessPolicyCreateUpdate,
		Read:   resourceIotHubDPSSharedAccessPolicyRead,
		Update: resourceIotHubDPSSharedAccessPolicyCreateUpdate,
		Delete: resourceIotHubDPSSharedAccessPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := iotdpsresource.ParseKeyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: iothubDpsSharedAccessPolicySchema(),
	}
}
