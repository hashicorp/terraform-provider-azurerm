// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumequotarules"
)

func TestValidateNetAppVolumeQuotaRules(t *testing.T) {
	cases := []struct {
		Name                 string
		VolumeQuotaRulesData volumequotarules.VolumeQuotaRulesProperties
		Errors               int
	}{
		{
			Name: "ValidateCorrectWindowsUserIDLong",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("S-1-5-21-4128827716-1023963696-1645503205-500"),
			},
			Errors: 0,
		},
		{
			Name: "ValidateCorrectWindowsGroupIDLong",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("S-1-5-21-4128827716-1023963696-1645503205-1024"),
			},
			Errors: 0,
		},
		{
			Name: "ValidateCorrectWindowsGroupIDShort",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("S-1-5-32-583"),
			},
			Errors: 0,
		},
		{
			Name: "ValidateCorrectWindowsGroupIDExtraShort",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("S-1-5-18"),
			},
			Errors: 0,
		},
		{
			Name: "ValidateIncorrectWindowsUserIDLongFails",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("S-1-5-21-4128827716-1023963696-1645503205-5007777777777777777"),
			},
			Errors: 1,
		},
		{
			Name: "ValidateIncorrectWindowsGroupIDLongFails",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("AA1-5-21-4128827716-1023963696-1645503205-1024"),
			},
			Errors: 1,
		},
		{
			Name: "ValidateIncorrectWindowsGroupIDShortFails",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("S-B-5-32-583"),
			},
			Errors: 1,
		},
		{
			Name: "ValidateIncorrectWindowsGroupIDExtraShortFails",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("A-1-5-18"),
			},
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := ValidateWindowsSID(*tc.VolumeQuotaRulesData.QuotaTarget, pointer.From(tc.VolumeQuotaRulesData.QuotaTarget))

			if len(errors) != tc.Errors {
				t.Fatalf("expected ValidateWindowsSID to return %d error(s) not %d\nError List: \n%v", tc.Errors, len(errors), errors)
			}
		})
	}
}
