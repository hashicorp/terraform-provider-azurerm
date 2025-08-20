// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
)

func FlattenLongTermBackUpScheduleDetails(longTermBackUpScheduleDetails *autonomousdatabases.LongTermBackUpScheduleDetails) []LongTermBackUpScheduleDetails {
	output := make([]LongTermBackUpScheduleDetails, 0)
	if longTermBackUpScheduleDetails != nil {
		return append(output, LongTermBackUpScheduleDetails{
			RepeatCadence:         string(pointer.From(longTermBackUpScheduleDetails.RepeatCadence)),
			TimeOfBackup:          pointer.From(longTermBackUpScheduleDetails.TimeOfBackup),
			RetentionPeriodInDays: pointer.From(longTermBackUpScheduleDetails.RetentionPeriodInDays),
			Enabled:               !pointer.From(longTermBackUpScheduleDetails.IsDisabled),
		})
	}
	return output
}

func expandCloneCustomerContacts(input []string) []autonomousdatabases.CustomerContact {
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

func flattenCloneCustomerContacts(input []autonomousdatabases.CustomerContact) []string {
	if len(input) == 0 {
		return nil
	}

	emails := make([]string, 0, len(input))
	for _, contact := range input {
		if contact.Email != "" {
			emails = append(emails, contact.Email)
		}
	}
	return emails
}
