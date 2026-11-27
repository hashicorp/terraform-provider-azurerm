// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_attached_network

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Dev Center Attached Network.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"dev_center_id": commonschema.ResourceIDReferenceRequiredForceNew(&devcenters.DevCenterId{}, "Dev Center"),

		"network_connection_id": commonschema.ResourceIDReferenceRequiredForceNew(&networkconnections.NetworkConnectionId{}, "Dev Center Network Connection"),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
