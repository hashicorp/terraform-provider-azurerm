// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataShareSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ShareName(),
		},

		"account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: account.ValidateAccountID,
		},

		"kind": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(share.PossibleValuesForShareKind(), false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"snapshot_schedule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.SnapshotScheduleName(),
					},

					"recurrence": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(synchronizationsetting.PossibleValuesForRecurrenceInterval(), false),
					},

					"start_time": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.IsRFC3339Time,
						DiffSuppressFunc: suppress.RFC3339Time,
					},
				},
			},
		},

		"terms": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}
