// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_subnet_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/associationsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (t Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApplicationLoadBalancerSubnetAssociationName(),
		},

		"application_load_balancer_id": commonschema.ResourceIDReferenceRequiredForceNew(&associationsinterface.TrafficControllerId{}, "Application Gateway for Containers"),

		"subnet_id": commonschema.ResourceIDReferenceRequired(&commonids.SubnetId{}, "Subnet"),

		"tags": commonschema.Tags(),
	}
}

func (t Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
