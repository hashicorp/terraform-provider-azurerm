// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
)

func ValidateNetAppVolumeGroupExportPolicyRule(rule volumegroups.ExportPolicyRule, protocolType string) []error {
	errors := make([]error, 0)

	// Validating that nfsv3 and nfsv4.1 are not enabled in the same rule
	if pointer.From(rule.Nfsv3) && pointer.From(rule.Nfsv41) {
		errors = append(errors, fmt.Errorf("'nfsv3 and nfsv4.1 cannot be enabled at the same time'"))
	}

	// Validating that nfsv3 and nfsv4.1 are not disabled in the same rule
	if !pointer.From(rule.Nfsv3) && !pointer.From(rule.Nfsv41) {
		errors = append(errors, fmt.Errorf("'nfsv3 and nfsv4.1 cannot be enabled at the same time'"))
	}

	// Validating that nfsv4.1 export policy is not set on nfsv3 volume
	if pointer.From(rule.Nfsv41) && strings.EqualFold(protocolType, string(ProtocolTypeNfsV3)) {
		errors = append(errors, fmt.Errorf("'nfsv4.1 export policy cannot be enabled on nfsv3 volume'"))
	}

	// Validating that nfsv3 export policy is not set on nfsv4.1 volume
	if pointer.From(rule.Nfsv3) && strings.EqualFold(protocolType, string(ProtocolTypeNfsV41)) {
		errors = append(errors, fmt.Errorf("'nfsv3 export policy cannot be enabled on nfsv4.1 volume'"))
	}

	return errors
}
