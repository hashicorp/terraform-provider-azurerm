// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ApiCenterServiceResourceModel struct {
	Name          string                                     `tfschema:"name"`
	ResourceGroup string                                     `tfschema:"resource_group_name"`
	Location      string                                     `tfschema:"location"`
	SkuName       string                                     `tfschema:"sku_name"`
	Identity      []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags          map[string]string                          `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &ApiCenterServiceResourceModel{}
}
