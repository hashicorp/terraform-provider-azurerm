// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumequotarules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func ValidateNetAppVolumeQuotaRule(ctx context.Context, volumeID volumes.VolumeId, client *clients.Client, rule *volumequotarules.VolumeQuotaRule) []error {
	errors := make([]error, 0)

	// Validating quota type matches volume type
	volumeClient := client.NetApp.VolumeClient

	volume, err := volumeClient.Get(ctx, volumeID)
	if err != nil && volume.HttpResponse.StatusCode != http.StatusNotFound {
		errors = append(errors, fmt.Errorf("'volume %v required for quota %v not found'", volumeID.VolumeName, rule.Name))
	}

	// NFS Volume must not have Windows SIDs
	if len(pointer.From(volume.Model.Properties.ProtocolTypes)) == 1 &&
		(findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "nfsv3") || findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "nfsv4.1")) {
		_, errList := ValidateWindowsSID(rule.Properties.QuotaTarget, pointer.From(rule.Properties.QuotaTarget))
		if len(errList) == 0 {
			errors = append(errors, fmt.Errorf("'nfs volume %v cannot have windows sid as quota targets, defined quota target is %v'", volumeID.VolumeName, rule.Properties.QuotaTarget))
		}
	}

	// CIFS Volume must not have Windows SIDs
	if len(pointer.From(volume.Model.Properties.ProtocolTypes)) == 1 && findStringInSlice(pointer.From(volume.Model.Properties.ProtocolTypes), "cifs") {
		_, errList := ValidateUnixUserIDOrGroupID(rule.Properties.QuotaTarget, pointer.From(rule.Properties.QuotaTarget))
		if len(errList) == 0 {
			errors = append(errors, fmt.Errorf("'cifs volume %v cannot have unix id as quota target, defined quota target is %v'", volumeID.VolumeName, rule.Properties.QuotaTarget))
		}
	}

	// DefaultUserQuota quota type cannot have quota target defined
	if pointer.From(rule.Properties.QuotaType) == volumequotarules.TypeDefaultUserQuota && rule.Properties.QuotaTarget != nil {
		errors = append(errors, fmt.Errorf("'defaultUserQuota cannot have quota target defined, defined quota target is %v'", rule.Properties.QuotaTarget))
	}

	// DefaultGroupQuota quota type cannot have quota target defined
	if pointer.From(rule.Properties.QuotaType) == volumequotarules.TypeDefaultGroupQuota && rule.Properties.QuotaTarget != nil {
		errors = append(errors, fmt.Errorf("'defaultGroupQuota cannot have quota target defined, defined quota target is %v'", rule.Properties.QuotaTarget))
	}

	return errors
}
