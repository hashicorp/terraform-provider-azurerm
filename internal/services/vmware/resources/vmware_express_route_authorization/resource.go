// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vmware_express_route_authorization

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/authorizations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVmwareExpressRouteAuthorizationCreate,
		Read:   resourceVmwareExpressRouteAuthorizationRead,
		Delete: resourceVmwareExpressRouteAuthorizationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := authorizations.ParseAuthorizationID(id)
			return err
		}),

		Schema: vmwareExpressRouteAuthorizationSchema(),
	}
}
