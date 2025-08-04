// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestValidateNetAppVolumeGroupExportPolicyRuleOracle(t *testing.T) {
	cases := []struct {
		Name     string
		Protocol string
		Rule     volumegroups.ExportPolicyRule
		Errors   int
	}{
		{
			Name:     "ValidateNFSv41EnabledOnNFSv41Volume",
			Protocol: string(ProtocolTypeNfsV41),
			Rule: volumegroups.ExportPolicyRule{
				Nfsv3:  pointer.To(false),
				Nfsv41: pointer.To(true),
			},
			Errors: 0,
		},
		{
			Name:     "ValidateNFSv3EnabledOnNFSv3Volume",
			Protocol: string(ProtocolTypeNfsV3),
			Rule: volumegroups.ExportPolicyRule{
				Nfsv3:  pointer.To(true),
				Nfsv41: utils.Bool(false),
			},
			Errors: 0,
		},
		{
			Name:     "ValidateBothProtocolsNotEnabledAtSameTimeOnNFSv41Volume",
			Protocol: string(ProtocolTypeNfsV41),
			Rule: volumegroups.ExportPolicyRule{
				Nfsv3:  pointer.To(true),
				Nfsv41: utils.Bool(true),
			},
			Errors: 2,
		},
		{
			Name:     "ValidateBothProtocolsNotEnabledAtSameTimeOnNFSv3Volume",
			Protocol: string(ProtocolTypeNfsV3),
			Rule: volumegroups.ExportPolicyRule{
				Nfsv3:  pointer.To(true),
				Nfsv41: utils.Bool(true),
			},
			Errors: 2,
		},
		{
			Name:     "ValidateBothProtocolsNotDisabledAtSameTimeOnNFSv3Volume",
			Protocol: string(ProtocolTypeNfsV3),
			Rule: volumegroups.ExportPolicyRule{
				Nfsv3:  pointer.To(false),
				Nfsv41: utils.Bool(false),
			},
			Errors: 1,
		},
		{
			Name:     "ValidateBothProtocolsNotDisabledAtSameTimeOnNFSv41Volume",
			Protocol: string(ProtocolTypeNfsV41),
			Rule: volumegroups.ExportPolicyRule{
				Nfsv3:  pointer.To(false),
				Nfsv41: utils.Bool(false),
			},
			Errors: 1,
		},
		{
			Name:     "ValidateNFSv3NotEnabledOnNFSv41Volume",
			Protocol: string(ProtocolTypeNfsV41),
			Rule: volumegroups.ExportPolicyRule{
				Nfsv3:  pointer.To(true),
				Nfsv41: utils.Bool(false),
			},
			Errors: 1,
		},
		{
			Name:     "ValidateNFSv41NotEnabledOnNFSv3Volume",
			Protocol: string(ProtocolTypeNfsV3),
			Rule: volumegroups.ExportPolicyRule{
				Nfsv3:  pointer.To(false),
				Nfsv41: utils.Bool(true),
			},
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			errors := ValidateNetAppVolumeGroupExportPolicyRule(tc.Rule, tc.Protocol)

			if len(errors) != tc.Errors {
				t.Fatalf("expected ValidateNetAppVolumeGroupOracleVolumes to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
