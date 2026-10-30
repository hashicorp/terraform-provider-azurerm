// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_quota_rule

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func (r Resource) ModelObject() any {
	return &models.NetAppVolumeQuotaRuleModel{}
}
