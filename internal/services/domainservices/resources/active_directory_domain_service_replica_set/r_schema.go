// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_replica_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func activeDirectoryDomainServiceReplicaSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"domain_service_id": {
			Description:  "The ID of the Domain Service for which to create this Replica Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DomainServiceID,
		},

		"location": commonschema.Location(),

		"subnet_id": {
			Description:  "The ID of the subnet in which to place this Replica Set.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"domain_controller_ip_addresses": {
			Description: "A list of subnet IP addresses for the domain controllers in this Replica Set, typically two.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"external_access_ip_address": {
			Description: "The publicly routable IP address for the domain controllers in this Replica Set.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"service_status": {
			Description: "The current service status for the replica set.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
