// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for the Microsoft SQL Virtual Machine Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 15),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sql_image_offer": {
			Description:  "The offer type of the marketplace image cluster to be used by the SQL Virtual Machine Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SqlImageOfferName,
		},

		"sql_image_sku": {
			Description:  "The SKU type of the marketplace image cluster to be used by the SQL Virtual Machine Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(sqlvirtualmachinegroups.PossibleValuesForSqlVMGroupImageSku(), false),
		},

		"wsfc_domain_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cluster_subnet_type": {
						Description:  "The subnet type of the SQL Virtual Machine cluster.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(sqlvirtualmachinegroups.PossibleValuesForClusterSubnetType(), false),
					},

					"fqdn": {
						Description:  "The fully qualified name of the domain.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"organizational_unit_path": {
						Description:  "The organizational Unit path in which the nodes and cluster will be present.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"cluster_bootstrap_account_name": {
						Description:  "The account name used for creating cluster.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"cluster_operator_account_name": {
						Description:  "The account name used for operating cluster.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sql_service_account_name": {
						Description:  "The account name under which SQL service will run on all participating SQL Virtual Machines in the cluster.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_account_url": {
						Description:  "The SAS URL to the Storage Container of the witness Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_account_primary_key": {
						Description:  "The primary key of the Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
