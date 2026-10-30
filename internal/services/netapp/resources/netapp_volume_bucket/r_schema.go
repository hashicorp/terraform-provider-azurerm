// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return helpers.NetAppBucketResourceCommonArguments()
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return helpers.NetAppBucketResourceCommonAttributes()
}
