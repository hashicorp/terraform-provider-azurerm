// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// ValidateNetAppVolumeProtocolConversion validates protocol conversion requirements
func ValidateNetAppVolumeProtocolConversion(oldProtocols, newProtocols []string, kerberosEnabled bool, dataReplication []interface{}, exportPolicyRules []interface{}) []error {
	var errors []error

	// Only validate if this is a protocol change and not initial creation
	if len(oldProtocols) == 0 || len(newProtocols) == 0 {
		return errors
	}

	// Check if old is NFSv3 and new is NFSv4.1 or vice versa
	oldHasNFSv3 := utils.SliceContainsValue(oldProtocols, "NFSv3")
	oldHasNFSv41 := utils.SliceContainsValue(oldProtocols, "NFSv4.1")
	newHasNFSv3 := utils.SliceContainsValue(newProtocols, "NFSv3")
	newHasNFSv41 := utils.SliceContainsValue(newProtocols, "NFSv4.1")

	// Only validate if this is an NFS protocol conversion
	isNFSProtocolChange := (oldHasNFSv3 && !oldHasNFSv41 && newHasNFSv41 && !newHasNFSv3) ||
		(oldHasNFSv41 && !oldHasNFSv3 && newHasNFSv3 && !newHasNFSv41)

	if !isNFSProtocolChange {
		return errors // No validation needed for non-NFS conversions
	}

	// Validate that destination volumes in cross-region replication cannot be converted
	for _, replication := range dataReplication {
		if replicationMap, ok := replication.(map[string]interface{}); ok {
			if endpointType, exists := replicationMap["endpoint_type"]; exists {
				if endpointTypeStr, ok := endpointType.(string); ok {
					if endpointTypeStr == "dst" {
						errors = append(errors, fmt.Errorf("cannot convert a destination volume in a cross-region replication relationship"))
					}
				}
			}
		}
	}

	// Validate Kerberos restriction for NFSv4.1 to NFSv3 conversion
	if oldHasNFSv41 && newHasNFSv3 && kerberosEnabled {
		errors = append(errors, fmt.Errorf("cannot convert an NFSv4.1 volume with Kerberos enabled to NFSv3"))
	}

	// Validate dual-protocol restriction
	if len(oldProtocols) > 1 || len(newProtocols) > 1 {
		errors = append(errors, fmt.Errorf("cannot change the NFS version of a dual-protocol volume"))
	}

	// Validate that this is not a dual-protocol conversion
	oldHasCIFS := utils.SliceContainsValue(oldProtocols, "CIFS")
	newHasCIFS := utils.SliceContainsValue(newProtocols, "CIFS")

	if oldHasCIFS || newHasCIFS {
		errors = append(errors, fmt.Errorf("cannot convert a single-protocol NFS volume to a dual-protocol volume, or the other way around"))
	}

	// During protocol conversion, export policies are typically updated in the same configuration
	// to match the new protocol, so we skip export policy validation during conversion
	// The Azure API will validate protocol compatibility at apply time

	return errors
}
