// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_availability_group_listener

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for the Microsoft SQL Virtual Machine Availability Group Listener.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 15),
		},

		"sql_virtual_machine_group_id": {
			Description:  "The ID of the SQL Virtual Machine Group to create the listener.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: availabilitygrouplisteners.ValidateSqlVirtualMachineGroupID,
		},

		"availability_group_name": {
			Description: "The name of the Availability Group.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"port": {
			Description:  "The port of the listener.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsPortNumber,
		},

		"load_balancer_configuration": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ExactlyOneOf: []string{"load_balancer_configuration", "multi_subnet_ip_configuration"},
			ForceNew:     true,
			MaxItems:     1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"load_balancer_id": {
						Description:  "The ID of the Load Balancer.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: loadbalancers.ValidateLoadBalancerID,
					},

					"private_ip_address": {
						Description:  "The private IP Address of the listener.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsIPAddress,
					},

					"probe_port": {
						Description:  "The probe port of the listener.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsPortNumber,
					},

					"sql_virtual_machine_ids": {
						Description: "A list of SQL Virtual Machine IDs.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: sqlvirtualmachines.ValidateSqlVirtualMachineID,
						},
					},

					"subnet_id": {
						Description:  "The ID of the Subnet to create the listener.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.SubnetID,
					},
				},
			},
		},

		"multi_subnet_ip_configuration": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			ExactlyOneOf: []string{"load_balancer_configuration", "multi_subnet_ip_configuration"},
			ForceNew:     true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"private_ip_address": {
						Description:  "The private IP Address of the listener.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsIPAddress,
					},

					"sql_virtual_machine_id": {
						Description:  "The ID of the SQL Virtual Machine.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: sqlvirtualmachines.ValidateSqlVirtualMachineID,
					},

					"subnet_id": {
						Description:  "The ID of the Subnet to create the listener.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.SubnetID,
					},
				},
			},
		},

		"replica": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sql_virtual_machine_id": {
						Description:  "The ID of the SQL Virtual Machine.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.AsGeneratedID(sqlvirtualmachines.ParseSqlVirtualMachineIDInsensitively),
					},

					"role": {
						Description:  "The replica role for the availability group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(availabilitygrouplisteners.PossibleValuesForRole(), false),
					},

					"commit": {
						Description:  "The replica commit mode for the availability group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(availabilitygrouplisteners.PossibleValuesForCommit(), false),
					},

					"failover_mode": {
						Description:  "The replica failover mode for the availability group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(availabilitygrouplisteners.PossibleValuesForFailover(), false),
					},

					"readable_secondary": {
						Description:  "The replica readable secondary mode for the availability group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(availabilitygrouplisteners.PossibleValuesForReadableSecondary(), false),
					},
				},
			},
			Set: ReplicaSchemaMsSqlVirtualMachineAvailabilityGroupListenerHash,
		},
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
