// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
)

type VolumeSpecNameSAPHana string

const (
	VolumeSpecNameSAPHanaData       VolumeSpecNameSAPHana = "data"
	VolumeSpecNameSAPHanaLog        VolumeSpecNameSAPHana = "log"
	VolumeSpecNameSAPHanaShared     VolumeSpecNameSAPHana = "shared"
	VolumeSpecNameSAPHanaDataBackup VolumeSpecNameSAPHana = "data-backup"
	VolumeSpecNameSAPHanaLogBackup  VolumeSpecNameSAPHana = "log-backup"
)

func PossibleValuesForVolumeSpecNameSAPHana() []string {
	return []string{
		string(VolumeSpecNameSAPHanaData),
		string(VolumeSpecNameSAPHanaLog),
		string(VolumeSpecNameSAPHanaShared),
		string(VolumeSpecNameSAPHanaDataBackup),
		string(VolumeSpecNameSAPHanaLogBackup),
	}
}

func RequiredVolumesForSAPHANA() []string {
	return []string{
		string(VolumeSpecNameSAPHanaData),
		string(VolumeSpecNameSAPHanaLog),
	}
}

func PossibleValuesForProtocolTypeVolumeGroupSAPHana() []string {
	return []string{
		string(ProtocolTypeNfsV41),
		string(ProtocolTypeNfsV3),
	}
}

func ValidateNetAppVolumeGroupSAPHanaVolumes(volumeList *[]volumegroups.VolumeGroupVolumeProperties) []error {
	errors := make([]error, 0)
	volumeSpecRepeatCount := make(map[string]int)
	applicationType := string(volumegroups.ApplicationTypeSAPNegativeHANA)

	// Validating each volume
	for _, volume := range pointer.From(volumeList) {
		// Get protocol list
		protocolTypeList := pointer.From(volume.Properties.ProtocolTypes)

		// Validate protocol list is not empty
		if len(protocolTypeList) == 0 {
			errors = append(errors, fmt.Errorf("'protocol type list cannot be empty'"))
		}

		for _, protocolType := range protocolTypeList {
			// Validate protocol list does not contain invalid protocols
			for _, protocol := range protocolTypeList {
				if !findStringInSlice(PossibleValuesForProtocolType(), protocolType) {
					errors = append(errors, fmt.Errorf("'protocol %v is invalid'", protocol))
				}
			}

			// Validate that protocol is valid for SAP Hana
			if !findStringInSlice(PossibleValuesForProtocolTypeVolumeGroupSAPHana(), protocolType) {
				errors = append(errors, fmt.Errorf("'protocol `%v` is invalid for SAP Hana'", protocolType))
			}

			// Can't be nfsv3 on data, log and share volumes
			if strings.EqualFold(protocolType, string(ProtocolTypeNfsV3)) &&
				(strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaData)) ||
					strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaShared)) ||
					strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaLog))) {
				errors = append(errors, fmt.Errorf("'nfsv3 on data, log and shared volumes for %v is not supported on volume %v'", applicationType, pointer.From(volume.Name)))
			}

			// Validating export policies
			if volume.Properties.ExportPolicy != nil {
				for _, rule := range pointer.From(volume.Properties.ExportPolicy.Rules) {
					errors = append(errors, ValidateNetAppVolumeGroupExportPolicyRule(rule, protocolType)...)
				}
			}
		}

		// Checking CRR rule that log cannot be DataProtection type
		if strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaLog)) &&
			volume.Properties.DataProtection != nil &&
			volume.Properties.DataProtection.Replication != nil &&
			strings.EqualFold(string(pointer.From(volume.Properties.DataProtection.Replication.EndpointType)), string(volumegroups.EndpointTypeDst)) {
			errors = append(errors, fmt.Errorf("'log volume spec type cannot be DataProtection type for %v on volume %v'", applicationType, pointer.From(volume.Name)))
		}

		// Validating that snapshot policies are not being created in a data protection volume
		if volume.Properties.DataProtection != nil &&
			volume.Properties.DataProtection.Snapshot != nil &&
			(volume.Properties.DataProtection.Replication != nil && strings.EqualFold(string(pointer.From(volume.Properties.DataProtection.Replication.EndpointType)), string(volumegroups.EndpointTypeDst))) {
			errors = append(errors, fmt.Errorf("'snapshot policy cannot be enabled on a data protection volume for %v on volume %v'", applicationType, pointer.From(volume.Name)))
		}

		// Validating that data-backup and log-backup don't have PPG defined
		if (strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaDataBackup)) ||
			strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaLogBackup))) &&
			pointer.From(volume.Properties.ProximityPlacementGroup) != "" {
			errors = append(errors, fmt.Errorf("'%v volume spec type cannot have PPG defined for %v on volume %v'", pointer.From(volume.Properties.VolumeSpecName), applicationType, pointer.From(volume.Name)))
		}

		// Validating that data, log and shared have PPG defined.
		if (strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaData)) ||
			strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaLog)) ||
			strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameSAPHanaShared))) &&
			pointer.From(volume.Properties.ProximityPlacementGroup) == "" {
			errors = append(errors, fmt.Errorf("'%v volume spec type must have PPG defined for %v on volume %v'", pointer.From(volume.Properties.VolumeSpecName), applicationType, pointer.From(volume.Name)))
		}

		// Adding volume spec name to hashmap for post volume loop check
		volumeSpecRepeatCount[pointer.From(volume.Properties.VolumeSpecName)] += 1
	}

	// Validating required volume spec types
	for _, requiredVolumeSpec := range RequiredVolumesForSAPHANA() {
		if _, ok := volumeSpecRepeatCount[requiredVolumeSpec]; !ok {
			errors = append(errors, fmt.Errorf("'required volume spec type %v is not present for %v'", requiredVolumeSpec, applicationType))
		}
	}

	// Validating that volume spec does not repeat
	for volumeSpecName, count := range volumeSpecRepeatCount {
		if count > 1 {
			errors = append(errors, fmt.Errorf("'volume spec type %v cannot be repeated for %v'", volumeSpecName, applicationType))
		}
	}

	return errors
}
