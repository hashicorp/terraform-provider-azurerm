// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumequotarules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func ValidateNetAppVolumeQuotaRule(ctx context.Context, volumeID volumes.VolumeId, client *clients.Client, rule *netAppModels.NetAppVolumeQuotaRuleModel) []error {
	errors := make([]error, 0)

	// Validating quota type matches volume type
	volumeClient := client.NetApp.VolumeClient

	volume, err := volumeClient.Get(ctx, volumeID)
	if err != nil {
		errors = append(errors, fmt.Errorf("retrieving %s: %v", volumeID, err))
		return errors
	}

	if volume.Model == nil || volume.HttpResponse.StatusCode == http.StatusNotFound {
		errors = append(errors, fmt.Errorf("volume %s not found", volumeID))
		return errors
	}

	// NFS Volume must not have Windows SIDs
	if len(pointer.From(volume.Model.Properties.ProtocolTypes)) == 1 &&
		(findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "nfsv3") || findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "nfsv4.1")) {
		if pointer.To(rule.QuotaTarget) != nil && rule.QuotaTarget != "" {
			_, errList := ValidateWindowsSID(&rule.QuotaTarget, rule.QuotaTarget)
			if len(errList) == 0 {
				errors = append(errors, fmt.Errorf("'nfs volume %v cannot have windows sid as quota targets, defined quota target is %v'", volumeID.VolumeName, rule.QuotaTarget))
			}
		}
	}

	// CIFS Volume must not have Windows SIDs
	if len(pointer.From(volume.Model.Properties.ProtocolTypes)) == 1 && findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "cifs") {
		if pointer.To(rule.QuotaTarget) != nil && rule.QuotaTarget != "" {
			_, errList := ValidateUnixUserIDOrGroupID(&rule.QuotaTarget, rule.QuotaTarget)
			if len(errList) == 0 {
				errors = append(errors, fmt.Errorf("'cifs volume %v cannot have unix id as quota target, defined quota target is %v'", volumeID.VolumeName, rule.QuotaTarget))
			}
		}
	}

	// Dual protocol volumes does not support group quotas
	if len(pointer.From(volume.Model.Properties.ProtocolTypes)) == 2 {
		if pointer.To(rule.QuotaType) != nil && (volumequotarules.Type(rule.QuotaType) == volumequotarules.TypeIndividualGroupQuota || volumequotarules.Type(rule.QuotaType) == volumequotarules.TypeDefaultGroupQuota) {
			errors = append(errors, fmt.Errorf("'dual protocol volume %v cannot have group quotas'", volumeID.VolumeName))
		}
	}

	// CIFS protocol volumes does not support group quotas
	if findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "cifs") {
		if pointer.To(rule.QuotaType) != nil && (volumequotarules.Type(rule.QuotaType) == volumequotarules.TypeIndividualGroupQuota || volumequotarules.Type(rule.QuotaType) == volumequotarules.TypeDefaultGroupQuota) {
			errors = append(errors, fmt.Errorf("'cifs volume %v cannot have group quotas'", volumeID.VolumeName))
		}
	}

	// Quota types and targets validations
	errors = append(errors, ValidateNetAppVolumeQuotaRuleQuotaType(pointer.To(volumequotarules.Type(rule.QuotaType)), pointer.To(rule.QuotaTarget))...)

	return errors
}

func ValidateNetAppVolumeQuotaRuleQuotaType(quotaType *volumequotarules.Type, quotaTarget *string) []error {
	errors := make([]error, 0)

	// Quota Target must be defined for IndividualUserQuota
	if quotaType != nil && pointer.From(quotaType) == volumequotarules.TypeIndividualUserQuota && (quotaTarget == nil || pointer.From(quotaTarget) == "") {
		errors = append(errors, fmt.Errorf("'individualUserQuota must have quota target defined'"))
	}

	// Quota Target must be defined for IndividualGroupQuota
	if quotaType != nil && pointer.From(quotaType) == volumequotarules.TypeIndividualGroupQuota && (quotaTarget == nil || pointer.From(quotaTarget) == "") {
		errors = append(errors, fmt.Errorf("'individualGroupQuota must have quota target defined'"))
	}

	// DefaultUserQuota quota type cannot have quota target defined
	if quotaType != nil && pointer.From(quotaType) == volumequotarules.TypeDefaultUserQuota && (quotaTarget != nil && pointer.From(quotaTarget) != "") {
		errors = append(errors, fmt.Errorf("'defaultUserQuota cannot have quota target defined, defined quota target is %v'", pointer.From(quotaTarget)))
	}

	// DefaultGroupQuota quota type cannot have quota target defined
	if quotaType != nil && pointer.From(quotaType) == volumequotarules.TypeDefaultGroupQuota && (quotaTarget != nil && pointer.From(quotaTarget) != "") {
		errors = append(errors, fmt.Errorf("'defaultGroupQuota cannot have quota target defined, defined quota target is %v'", pointer.From(quotaTarget)))
	}

	return errors
}
