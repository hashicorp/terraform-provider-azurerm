// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_center_network_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Dev Center Network Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevCenterNetworkConnectionName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"domain_join_type": {
			Description:  "The Azure Active Directory Join type.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(networkconnections.PossibleValuesForDomainJoinType(), false),
		},

		"subnet_id": commonschema.ResourceIDReferenceRequired(&commonids.SubnetId{}, "Subnet"),

		"domain_name": {
			Description:  "The name of the Azure Active Directory domain.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.DevCenterNetworkConnectionDomainName,
		},

		"domain_password": {
			Description:  "The password for the account used to join domain.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"domain_username": {
			Description:  "The username of the Azure Active Directory account (user or service account) that has permissions to create computer objects in Active Directory.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.DevCenterNetworkConnectionDomainUsername,
		},

		"organization_unit": {
			Description:  "The Azure Active Directory domain Organization Unit (OU).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
