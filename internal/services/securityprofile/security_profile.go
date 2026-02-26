// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securityprofile

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

// Values captures the Terraform view of a security profile.
type Values struct {
	HostEncryption *bool
	SecurityType   *string
	SecureBoot     *bool
	VTPM           *bool
}

// Accessor exposes the getters and setters required to read/write security profile fields on service-specific types.
// FromBlock converts the schema block into a Values struct and reports if the block was configured.
func FromBlock(blocks []interface{}) *Values {
	if len(blocks) == 0 {
		return nil
	}

	item := blocks[0].(map[string]interface{})
	var securityProfile Values

	if v, ok := item["host_encryption_enabled"]; ok {
		securityProfile.HostEncryption = pointer.To(v.(bool))
	}

	if v, ok := item["security_type"]; ok {
		securityProfile.SecurityType = pointer.To(v.(string))
	}

	if v, ok := item["secure_boot_enabled"]; ok {
		securityProfile.SecureBoot = pointer.To(v.(bool))
	}

	if v, ok := item["vtpm_enabled"]; ok {
		securityProfile.VTPM = pointer.To(v.(bool))
	}

	return &securityProfile
}

// ToBlock flattens Values back into the Terraform schema representation.
func ToBlock(values Values) []interface{} {
	securityProfile := make([]interface{}, 0)
	securityProfile = append(securityProfile, map[string]interface{}{
		"host_encryption_enabled": pointer.From(values.HostEncryption),
		"security_type":           pointer.From(values.SecurityType),
		"secure_boot_enabled":     pointer.From(values.SecureBoot),
		"vtpm_enabled":            pointer.From(values.VTPM),
	})
	return securityProfile
}
