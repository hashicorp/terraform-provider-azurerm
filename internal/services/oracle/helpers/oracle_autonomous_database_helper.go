// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
)

func FlattenLongTermBackUpScheduleDetails(longTermBackUpScheduleDetails *autonomousdatabases.LongTermBackUpScheduleDetails) []LongTermBackUpScheduleDetails {
	output := make([]LongTermBackUpScheduleDetails, 0)
	if longTermBackUpScheduleDetails != nil {
		return append(output, LongTermBackUpScheduleDetails{
			RepeatCadence:         pointer.FromEnum(longTermBackUpScheduleDetails.RepeatCadence),
			TimeOfBackup:          pointer.From(longTermBackUpScheduleDetails.TimeOfBackup),
			RetentionPeriodInDays: pointer.From(longTermBackUpScheduleDetails.RetentionPeriodInDays),
			Enabled:               !pointer.From(longTermBackUpScheduleDetails.IsDisabled),
		})
	}
	return output
}

func ExpandCloneCustomerContacts(input []string) []autonomousdatabases.CustomerContact {
	if len(input) == 0 {
		return nil
	}

	contacts := make([]autonomousdatabases.CustomerContact, 0, len(input))
	for _, email := range input {
		contacts = append(contacts, autonomousdatabases.CustomerContact{
			Email: email,
		})
	}
	return contacts
}

func FlattenConnectionStrings(connStrings *autonomousdatabases.ConnectionStringType) []string {
	flattened := make([]string, 0)

	if connStrings == nil {
		return flattened
	}
	allConnStrings := connStrings.AllConnectionStrings
	if allConnStrings == nil {
		return flattened
	}

	if allConnStrings.High != nil {
		flattened = append(flattened, *allConnStrings.High)
	}

	if allConnStrings.Medium != nil {
		flattened = append(flattened, *allConnStrings.Medium)
	}

	if allConnStrings.Low != nil {
		flattened = append(flattened, *allConnStrings.Low)
	}

	return flattened
}
