// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azure

import (
	"fmt"
	"time"
)

// IsRetired determins if a the current clients locale date/time is before or after the passed date.
func IsRetired(year int, month time.Month, day int) (bool, error) {
	var retired bool

	// Set time zone location to PST/PDT (e.g., Redmond)...
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return false, fmt.Errorf("unable to load location 'America/Los_Angeles'")
	}

	input := time.Date(year, month, day, 0, 0, 0, 0, location)

	// Set client time to midnight UTC...
	clientNow := time.Now().UTC().Truncate(24 * time.Hour)
	clienttUTC := time.Date(clientNow.Year(), clientNow.Month(), clientNow.Day(), 0, 0, 0, 0, time.UTC)

	// Normalize the retiredDate from PST/PDT time zone to UTC...
	if clienttUTC.Equal(input.UTC()) || clienttUTC.After(input.UTC()) {
		retired = true
	}

	return retired, nil
}
