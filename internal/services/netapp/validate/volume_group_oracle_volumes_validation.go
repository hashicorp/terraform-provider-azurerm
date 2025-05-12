// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
)

type VolumeSpecNameOracle string

const (
	VolumeSpecNameOracleData1  VolumeSpecNameOracle = "ora-data1"
	VolumeSpecNameOracleData2  VolumeSpecNameOracle = "ora-data2"
	VolumeSpecNameOracleData3  VolumeSpecNameOracle = "ora-data3"
	VolumeSpecNameOracleData4  VolumeSpecNameOracle = "ora-data4"
	VolumeSpecNameOracleData5  VolumeSpecNameOracle = "ora-data5"
	VolumeSpecNameOracleData6  VolumeSpecNameOracle = "ora-data6"
	VolumeSpecNameOracleData7  VolumeSpecNameOracle = "ora-data7"
	VolumeSpecNameOracleData8  VolumeSpecNameOracle = "ora-data8"
	VolumeSpecNameOracleLog    VolumeSpecNameOracle = "ora-log"
	VolumeSpecNameOracleMirror VolumeSpecNameOracle = "ora-log-mirror"
	VolumeSpecNameOracleBinary VolumeSpecNameOracle = "ora-binary"
	VolumeSpecNameOracleBackup VolumeSpecNameOracle = "ora-backup"
)

func PossibleValuesForVolumeSpecNameOracle() []string {
	return []string{
		string(VolumeSpecNameOracleData1),
		string(VolumeSpecNameOracleData2),
		string(VolumeSpecNameOracleData3),
		string(VolumeSpecNameOracleData4),
		string(VolumeSpecNameOracleData5),
		string(VolumeSpecNameOracleData6),
		string(VolumeSpecNameOracleData7),
		string(VolumeSpecNameOracleData8),
		string(VolumeSpecNameOracleLog),
		string(VolumeSpecNameOracleMirror),
		string(VolumeSpecNameOracleBinary),
		string(VolumeSpecNameOracleBackup),
	}
}

func RequiredVolumesForOracle() []string {
	return []string{
		string(VolumeSpecNameOracleData1),
		string(VolumeSpecNameOracleLog),
	}
}

func PossibleValuesForProtocolTypeVolumeGroupOracle() []string {
	return []string{
		string(ProtocolTypeNfsV41),
		string(ProtocolTypeNfsV3),
	}
}

func ValidateNetAppVolumeGroupOracleVolumes(volumeList *[]volumegroups.VolumeGroupVolumeProperties) []error {
	errors := make([]error, 0)
	expectedZone := ""
	expectedPpgId := ""
	expectedKeyVaultPrivateEndpointResourceId := ""
	expectedEncryptionKeySource := ""
	volumeSpecRepeatCount := make(map[string]int)
	applicationType := string(volumegroups.ApplicationTypeORACLE)

	// Validating each volume
	for _, volume := range pointer.From(volumeList) {
		// Get protocol list
		protocolTypeList := pointer.From(volume.Properties.ProtocolTypes)

		// Validate protocol list is not empty
		if len(protocolTypeList) == 0 {
			errors = append(errors, fmt.Errorf("'protocol type list cannot be empty'"))
		}

		for _, protocol := range protocolTypeList {
			// Validate protocol list does not contain invalid protocols
			if !findStringInSlice(PossibleValuesForProtocolType(), protocol) {
				errors = append(errors, fmt.Errorf("'protocol %v is invalid'", protocol))
			}

			// Validate that protocol is valid for Oracle
			if !findStringInSlice(PossibleValuesForProtocolTypeVolumeGroupOracle(), protocol) {
				errors = append(errors, fmt.Errorf("'protocol %v is invalid for Oracle'", protocol))
			}

			// Validating export policies
			if volume.Properties.ExportPolicy != nil {
				for _, rule := range pointer.From(volume.Properties.ExportPolicy.Rules) {
					errors = append(errors, ValidateNetAppVolumeGroupExportPolicyRule(rule, protocol)...)
				}
			}
		}

		// Checking CRR rule that log cannot be DataProtection type
		if strings.EqualFold(pointer.From(volume.Properties.VolumeSpecName), string(VolumeSpecNameOracleLog)) &&
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

		// Validating that zone and proximity placement group are not specified together
		if (pointer.From(volume.Zones) != nil || len(pointer.From(volume.Zones)) > 0) && pointer.From(volume.Properties.ProximityPlacementGroup) != "" {
			errors = append(errors, fmt.Errorf("'zone and proximity_placement_group_id cannot be specified together on volume %v'", pointer.From(volume.Name)))
		}

		// Getting the first zone for validations, all volumes must be in the same zone
		if expectedZone == "" {
			if volume.Zones != nil && len(pointer.From(volume.Zones)) > 0 {
				expectedZone = pointer.From(volume.Zones)[0]
			}
		}

		// Validating that all volumes are in the same zone
		if volume.Zones != nil && len(pointer.From(volume.Zones)) > 0 && pointer.From(volume.Zones)[0] != expectedZone {
			errors = append(errors, fmt.Errorf("'zone must be the same on all volumes of this volume group, volume %v zone is %v'", pointer.From(volume.Name), pointer.From(volume.Zones)[0]))
		}

		// Getting the first PPG for validations, all volumes must be in the same PPG
		if expectedPpgId == "" {
			if volume.Properties.ProximityPlacementGroup != nil {
				expectedPpgId = pointer.From(volume.Properties.ProximityPlacementGroup)
			}
		}

		// Validating that all volumes are in the same PPG
		if volume.Properties.ProximityPlacementGroup != nil && pointer.From(volume.Properties.ProximityPlacementGroup) != expectedPpgId {
			errors = append(errors, fmt.Errorf("'proximity_placement_group_id must be the same on all volumes of this volume group, volume %v ppg id is %v'", pointer.From(volume.Name), pointer.From(volume.Properties.ProximityPlacementGroup)))
		}

		// Validating that encryption_key_source when key source is AKV has  key_vault_private_endpoint_id specified
		if pointer.From(volume.Properties.EncryptionKeySource) == volumegroups.EncryptionKeySourceMicrosoftPointKeyVault && pointer.From(volume.Properties.KeyVaultPrivateEndpointResourceId) == "" {
			errors = append(errors, fmt.Errorf("'encryption_key_source as microsoft.keyvault must have key_vault_private_endpoint_id specified on volume %v'", pointer.From(volume.Name)))
		}

		// Validating that encryption_key_source is set when key_vault_private_endpoint_id is specified
		if pointer.From(volume.Properties.KeyVaultPrivateEndpointResourceId) != "" && pointer.From(volume.Properties.EncryptionKeySource) == "" {
			errors = append(errors, fmt.Errorf("'encryption_key_source must be set when key_vault_private_endpoint_id is specified on volume %v'", pointer.From(volume.Name)))
		}

		// Getting the first KeyVaultPrivateEndpointResourceId for validations, all volumes must have the same KeyVaultPrivateEndpointResourceId
		if expectedKeyVaultPrivateEndpointResourceId == "" {
			if volume.Properties.KeyVaultPrivateEndpointResourceId != nil {
				expectedKeyVaultPrivateEndpointResourceId = pointer.From(volume.Properties.KeyVaultPrivateEndpointResourceId)
			}
		}

		// Validating that all volumes have the same KeyVaultPrivateEndpointResourceId
		if volume.Properties.KeyVaultPrivateEndpointResourceId != nil && pointer.From(volume.Properties.KeyVaultPrivateEndpointResourceId) != expectedKeyVaultPrivateEndpointResourceId {
			errors = append(errors, fmt.Errorf("'key_vault_private_endpoint_id must be the same on all volumes of this volume group, volume %v key vault private endpoint id is %v'", pointer.From(volume.Name), pointer.From(volume.Properties.KeyVaultPrivateEndpointResourceId)))
		}

		// Getting the first EncryptionKeySource for validations, all volumes must have the same EncryptionKeySource
		if expectedEncryptionKeySource == "" {
			if volume.Properties.EncryptionKeySource != nil {
				expectedEncryptionKeySource = string(pointer.From(volume.Properties.EncryptionKeySource))
			}
		}

		// Validating that all volumes have the same EncryptionKeySource
		if volume.Properties.EncryptionKeySource != nil && pointer.From(volume.Properties.EncryptionKeySource) != volumegroups.EncryptionKeySource(expectedEncryptionKeySource) {
			errors = append(errors, fmt.Errorf("'encryption_key_source must be the same on all volumes of this volume group, volume %v encryption key source is %v'", pointer.From(volume.Name), pointer.From(volume.Properties.EncryptionKeySource)))
		}

		// Validate that volume networkFeature is set to standard when volume.Properties.EncryptionKeySource == "Microsoft.KeyVault"
		if volume.Properties.EncryptionKeySource != nil && pointer.From(volume.Properties.EncryptionKeySource) == volumegroups.EncryptionKeySourceMicrosoftPointKeyVault && pointer.From(volume.Properties.NetworkFeatures) != volumegroups.NetworkFeaturesStandard {
			errors = append(errors, fmt.Errorf("'network_feature must be set to standard when encryption_key_source is set to Microsoft.KeyVault on volume %v, current value is %v'", pointer.From(volume.Name), pointer.From(volume.Properties.NetworkFeatures)))
		}

		// Adding volume spec name to hashmap for post volume loop check
		volumeSpecRepeatCount[pointer.From(volume.Properties.VolumeSpecName)] += 1
	}

	// Validating required volume spec types
	for _, requiredVolumeSpec := range RequiredVolumesForOracle() {
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
