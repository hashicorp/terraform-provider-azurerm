// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backuppolicy"
)

const (
	minDataDailyBackupsToKeep   int64 = 10
	minDataWeeklyBackupsToKeep  int64 = 10
	minDataMonthlyBackupsToKeep int64 = 10

	maxDataDailyBackupsToKeep   int64 = 1019
	maxDataWeeklyBackupsToKeep  int64 = 1019
	maxDataMonthlyBackupsToKeep int64 = 1019
)

func TestValidateNetAppBackupPolicyCombinedRetention(t *testing.T) {
	cases := []struct {
		Name             string
		BackupPolicyData backuppolicy.BackupPolicyProperties
		Errors           int
	}{
		{
			Name: "ValidateCombinedRetentionWithValidValues",
			BackupPolicyData: backuppolicy.BackupPolicyProperties{
				DailyBackupsToKeep:   pointer.To(minDataDailyBackupsToKeep),
				WeeklyBackupsToKeep:  pointer.To(minDataWeeklyBackupsToKeep),
				MonthlyBackupsToKeep: pointer.To(minDataMonthlyBackupsToKeep),
			},
			Errors: 0,
		},
		{
			Name: "ValidateCombinedRetentionWithValidValuesMaximumReachedFailure",
			BackupPolicyData: backuppolicy.BackupPolicyProperties{
				DailyBackupsToKeep:   pointer.To(maxDataDailyBackupsToKeep),
				WeeklyBackupsToKeep:  pointer.To(maxDataWeeklyBackupsToKeep),
				MonthlyBackupsToKeep: pointer.To(maxDataMonthlyBackupsToKeep),
			},
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			errors := ValidateNetAppBackupPolicyCombinedRetention(pointer.From(tc.BackupPolicyData.DailyBackupsToKeep), pointer.From(tc.BackupPolicyData.WeeklyBackupsToKeep), pointer.From(tc.BackupPolicyData.MonthlyBackupsToKeep))

			if len(errors) != tc.Errors {
				t.Fatalf("expected ValidateUnixUserIDOrGroupID to return %d error(s) not %d\nError List: \n%v", tc.Errors, len(errors), errors)
			}
		})
	}
}
