// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumequotarules"
)

func TestValidateNetAppVolumeQuotaRulesQuotaType(t *testing.T) {
	cases := []struct {
		Name                 string
		VolumeQuotaRulesData volumequotarules.VolumeQuotaRulesProperties
		Errors               int
	}{
		{
			Name: "ValidateIndividualUserQuotaTargetIsDefined",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("1001"),
				QuotaType:   pointer.To(volumequotarules.TypeIndividualUserQuota),
			},
			Errors: 0,
		},
		{
			Name: "ValidateIndividualGroupQuotaTargetIsDefined",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("1001"),
				QuotaType:   pointer.To(volumequotarules.TypeIndividualGroupQuota),
			},
			Errors: 0,
		},
		{
			Name: "ValidateDefaultUserQuotaTargetIsNotDefined",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaType: pointer.To(volumequotarules.TypeDefaultUserQuota),
			},
			Errors: 0,
		},
		{
			Name: "ValidateIndividualGroupQuotaTargetIsNotDefined",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaType: pointer.To(volumequotarules.TypeDefaultGroupQuota),
			},
			Errors: 0,
		},

		{
			Name: "ValidateIndividualUserQuotaTargetIsDefinedFailsWhenMissingTarget",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaType: pointer.To(volumequotarules.TypeIndividualUserQuota),
			},
			Errors: 1,
		},
		{
			Name: "ValidateIndividualGroupQuotaTargetIsDefinedFailsWhenMissingTarget",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaType: pointer.To(volumequotarules.TypeIndividualGroupQuota),
			},
			Errors: 1,
		},
		{
			Name: "ValidateDefaultUserQuotaTargetFailsWhenDefiningTarget",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("1001"),
				QuotaType:   pointer.To(volumequotarules.TypeDefaultUserQuota),
			},
			Errors: 1,
		},
		{
			Name: "ValidateIndividualGroupQuotaTargetFailsWhenDefiningTarget",
			VolumeQuotaRulesData: volumequotarules.VolumeQuotaRulesProperties{
				QuotaTarget: pointer.To("1001"),
				QuotaType:   pointer.To(volumequotarules.TypeDefaultGroupQuota),
			},
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			errors := ValidateNetAppVolumeQuotaRuleQuotaType(tc.VolumeQuotaRulesData.QuotaType, tc.VolumeQuotaRulesData.QuotaTarget)

			if len(errors) != tc.Errors {
				t.Fatalf("expected ValidateNetAppVolumeQuotaRuleQuotaType to return %d error(s) not %d\nError List: \n%v", tc.Errors, len(errors), errors)
			}
		})
	}
}
