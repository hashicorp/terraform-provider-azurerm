// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
)

const (
	maximumAllowedCombinedRetentionCount = 1019
)

func ValidateNetAppBackupPolicyCombinedRetention(dailyBackupsToKeep, weeklyBackupsToKeep, monthlyBackupsToKeep int64) []error {
	errors := make([]error, 0)

	// Validates that the combined retention count is less than the maximum allowed
	if dailyBackupsToKeep+weeklyBackupsToKeep+monthlyBackupsToKeep > maximumAllowedCombinedRetentionCount {
		errors = append(errors, fmt.Errorf("the combined retention count of daily_backups_to_keep, weekly_backups_to_keep, and monthly_backups_to_keep must be less than or equal to %d", maximumAllowedCombinedRetentionCount))
	}

	return errors
}

// TODO: Validating that the policy attached to a secondary CRR volume (destination) is not enabled for backup
func ValidateNetAppBackupPolicyForSecondaryCRRVolume(backupPolicyEnabled bool) []error {
	errors := make([]error, 0)

	if backupPolicyEnabled {
		errors = append(errors, fmt.Errorf("backup policy cannot be enabled on a secondary CRR volume"))
	}

	return errors
}
