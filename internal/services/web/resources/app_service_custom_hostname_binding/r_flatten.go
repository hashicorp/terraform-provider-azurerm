// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_custom_hostname_binding

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAppServiceCustomHostnameBindingFlatten(d *pluginsdk.ResourceData, id *webapps.HostNameBindingId, model *webapps.HostNameBinding) error {
	d.Set("hostname", id.HostNameBindingName)
	d.Set("app_service_name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("ssl_state", pointer.FromEnum(props.SslState))
			d.Set("thumbprint", props.Thumbprint)
			d.Set("virtual_ip", props.VirtualIP)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
