// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the NetApp Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PoolName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Description:  "The name of the NetApp account in which the NetApp Pool should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"service_level": {
			Description: "The service level of the file system.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(capacitypools.ServiceLevelPremium),
				string(capacitypools.ServiceLevelStandard),
				string(capacitypools.ServiceLevelUltra),
				string(capacitypools.ServiceLevelFlexible),
			}, false),
		},

		"size_in_tb": {
			Description:  "Provisioned size of the pool in TB.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 2048),
		},

		"qos_type": {
			Description:  "QoS Type of the pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(capacitypools.QosTypeAuto),
			ValidateFunc: validation.StringInSlice(capacitypools.PossibleValuesForQosType(), false),
		},

		"encryption_type": {
			Description:  "The encryption type of the pool.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      capacitypools.EncryptionTypeSingle,
			ValidateFunc: validation.StringInSlice(capacitypools.PossibleValuesForEncryptionType(), false),
		},

		"cool_access_enabled": {
			Description: "Whether the NetApp Pool can hold cool access enabled volumes.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"custom_throughput_mibps": {
			Description:  "The custom throughput for the pool in MiB/s. Minimum value is `128`. This field can only be set when `service_level` is set to `Flexible` and `qos_type` is set to `Manual`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(128),
		},

		"tags": commonschema.Tags(),
	}
}
