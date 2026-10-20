// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/capacitypools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func netappPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PoolName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"service_level": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(capacitypools.ServiceLevelPremium),
				string(capacitypools.ServiceLevelStandard),
				string(capacitypools.ServiceLevelUltra),
				string(capacitypools.ServiceLevelFlexible),
			}, false),
		},

		"size_in_tb": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 2048),
		},

		"qos_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(capacitypools.QosTypeAuto),
			ValidateFunc: validation.StringInSlice(capacitypools.PossibleValuesForQosType(), false),
		},

		"encryption_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      capacitypools.EncryptionTypeSingle,
			ValidateFunc: validation.StringInSlice(capacitypools.PossibleValuesForEncryptionType(), false),
		},

		"cool_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"custom_throughput_mibps": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(128),
		},

		"tags": commonschema.Tags(),
	}
}
