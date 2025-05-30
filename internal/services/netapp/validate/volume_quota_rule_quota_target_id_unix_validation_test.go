// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumequotarules"
)

func TestValidateNetAppVolumeQuotaRulesUnix(t *testing.T) {
	cases := []struct {
		Name                 string
		VolumeQuotaRulesData volumequotarules.VolumeQuotaRulesProperties
		Errors               int
	}{
		{
			Name: "ValidateCorrectUnixUserGroupID",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("1001"),
			},
			Errors: 0,
		},
		{
			Name: "ValidateUnixUserGroupIDWithLettersFail",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("a1001"),
			},
			Errors: 1,
		},
		{
			Name: "ValidateUnixUserGroupIDWithMoreThanMaxCharsFail",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("4294967296"),
			},
			Errors: 1,
		},
		{
			Name: "ValidateUnixUserGroupIDWithEmptyStringFail",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To(""),
			},
			Errors: 1,
		},
		{
			Name: "ValidateUnixUserGroupIDWithNegativeNumberFail",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("-1"),
			},
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := ValidateUnixUserIDOrGroupID(*tc.VolumeQuotaRulesData.QuotaTarget, pointer.From(tc.VolumeQuotaRulesData.QuotaTarget))

			if len(errors) != tc.Errors {
				t.Fatalf("expected ValidateUnixUserIDOrGroupID to return %d error(s) not %d\nError List: \n%v", tc.Errors, len(errors), errors)
			}
		})
	}
}
