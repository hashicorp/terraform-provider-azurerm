// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumequotarules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	netappClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func ValidateNetAppVolumeQuotaRule(ctx context.Context, volumeID volumes.VolumeId, client any, rule *models.NetAppVolumeQuotaRuleModel) []error {
	errs := make([]error, 0)

	// Validating quota type matches volume type
	volumeClient := netappClient.From(client).VolumeClient

	volume, err := volumeClient.Get(ctx, volumeID)
	if err != nil {
		errs = append(errs, fmt.Errorf("retrieving %s: %v", volumeID, err))
		return errs
	}

	if volume.Model == nil || volume.HttpResponse.StatusCode == http.StatusNotFound {
		errs = append(errs, fmt.Errorf("volume %s not found", volumeID))
		return errs
	}

	// NFS Volume must not have Windows SIDs
	if len(pointer.From(volume.Model.Properties.ProtocolTypes)) == 1 &&
		(findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "nfsv3") || findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "nfsv4.1")) {
		if pointer.To(rule.QuotaTarget) != nil && rule.QuotaTarget != "" {
			_, errList := ValidateWindowsSID(&rule.QuotaTarget, rule.QuotaTarget)
			if len(errList) == 0 {
				errs = append(errs, fmt.Errorf("'nfs volume %v cannot have windows sid as quota targets, defined quota target is %v'", volumeID.VolumeName, rule.QuotaTarget))
			}
		}
	}

	// CIFS Volume must not have Windows SIDs
	if len(pointer.From(volume.Model.Properties.ProtocolTypes)) == 1 && findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "cifs") {
		if pointer.To(rule.QuotaTarget) != nil && rule.QuotaTarget != "" {
			_, errList := ValidateUnixUserIDOrGroupID(&rule.QuotaTarget, rule.QuotaTarget)
			if len(errList) == 0 {
				errs = append(errs, fmt.Errorf("'cifs volume %v cannot have unix id as quota target, defined quota target is %v'", volumeID.VolumeName, rule.QuotaTarget))
			}
		}
	}

	// Dual protocol volumes does not support group quotas
	if len(pointer.From(volume.Model.Properties.ProtocolTypes)) == 2 {
		if pointer.To(rule.QuotaType) != nil && (volumequotarules.QuotaType(rule.QuotaType) == volumequotarules.QuotaTypeIndividualGroupQuota || volumequotarules.QuotaType(rule.QuotaType) == volumequotarules.QuotaTypeDefaultGroupQuota) {
			errs = append(errs, fmt.Errorf("'dual protocol volume %v cannot have group quotas'", volumeID.VolumeName))
		}
	}

	// CIFS protocol volumes does not support group quotas
	if findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "cifs") {
		if pointer.To(rule.QuotaType) != nil && (volumequotarules.QuotaType(rule.QuotaType) == volumequotarules.QuotaTypeIndividualGroupQuota || volumequotarules.QuotaType(rule.QuotaType) == volumequotarules.QuotaTypeDefaultGroupQuota) {
			errs = append(errs, fmt.Errorf("'cifs volume %v cannot have group quotas'", volumeID.VolumeName))
		}
	}

	// Quota types and targets validations
	errs = append(errs, ValidateNetAppVolumeQuotaRuleQuotaType(pointer.ToEnum[volumequotarules.QuotaType](rule.QuotaType), pointer.To(rule.QuotaTarget))...)

	return errs
}

func ValidateNetAppVolumeQuotaRuleQuotaType(quotaType *volumequotarules.QuotaType, quotaTarget *string) []error {
	errs := make([]error, 0)

	// Quota Target must be defined for IndividualUserQuota
	if quotaType != nil && pointer.From(quotaType) == volumequotarules.QuotaTypeIndividualUserQuota && (quotaTarget == nil || pointer.From(quotaTarget) == "") {
		errs = append(errs, errors.New("'individualUserQuota must have quota target defined'"))
	}

	// Quota Target must be defined for IndividualGroupQuota
	if quotaType != nil && pointer.From(quotaType) == volumequotarules.QuotaTypeIndividualGroupQuota && (quotaTarget == nil || pointer.From(quotaTarget) == "") {
		errs = append(errs, errors.New("'individualGroupQuota must have quota target defined'"))
	}

	// DefaultUserQuota quota type cannot have quota target defined
	if quotaType != nil && pointer.From(quotaType) == volumequotarules.QuotaTypeDefaultUserQuota && (quotaTarget != nil && pointer.From(quotaTarget) != "") {
		errs = append(errs, fmt.Errorf("'defaultUserQuota cannot have quota target defined, defined quota target is %v'", pointer.From(quotaTarget)))
	}

	// DefaultGroupQuota quota type cannot have quota target defined
	if quotaType != nil && pointer.From(quotaType) == volumequotarules.QuotaTypeDefaultGroupQuota && (quotaTarget != nil && pointer.From(quotaTarget) != "") {
		errs = append(errs, fmt.Errorf("'defaultGroupQuota cannot have quota target defined, defined quota target is %v'", pointer.From(quotaTarget)))
	}

	return errs
}
