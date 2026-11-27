// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_quota_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumequotarules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Volume Quota Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.VolumeQuotaRuleName,
		},

		"location": commonschema.Location(),

		"volume_id": {
			Description:  "The NetApp volume ID where the Volume Quota Rule is assigned to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: volumequotarules.ValidateVolumeID,
		},

		"quota_type": {
			Description:  "Quota type. Please note that `IndividualGroupQuota` and `DefaultGroupQuota` are not applicable to SMB and dual-protocol volumes.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(volumequotarules.PossibleValuesForQuotaType(), false),
		},

		"quota_size_in_kib": {
			Description:  "Quota size in kibibytes.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(4),
		},

		"quota_target": {
			Description:  "Quota Target. This can be Unix UID/GID for NFSv3/NFSv4.1 volumes and Windows User SID for CIFS based volumes.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
