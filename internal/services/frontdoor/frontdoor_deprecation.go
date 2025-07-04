// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"time"
)

var (
	CreateDeprecationMessage string = "the creation of new Frontdoor resources is no longer permitted following its deprecation on April 1, 2025. However, modifications to existing Frontdoor resources remain supported until the API reaches full retirement on March 31, 2027"
	FullyRetiredMessage      string = "as of March 31, 2027, Frontdoor resources have reached full retirement and are no longer supported by Azure"
)

// isAPIFieldDeprecated checks if an API field has been deprecated based on the deprecation date
// compared to the current time in the client's timezone versus Los Angeles timezone.
// deprecationDate should be in "2006-01-02" format (YYYY-MM-DD)
func isAPIFieldDeprecated(deprecationDate string) bool {
	// Parse the deprecation date
	deprecationTime, err := time.Parse("2006-01-02", deprecationDate)
	if err != nil {
		// If we can't parse the date, assume it's not deprecated to be safe
		return false
	}

	// Load Los Angeles timezone
	losAngelesLocation, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		// Fallback to UTC if we can't load LA timezone
		losAngelesLocation = time.UTC
	}

	// Get current time in client's local timezone
	clientNow := time.Now()

	// Get current time in Los Angeles timezone
	losAngelesNow := time.Now().In(losAngelesLocation)

	// Set deprecation time to start of day in Los Angeles timezone
	deprecationTimeLA := time.Date(
		deprecationTime.Year(),
		deprecationTime.Month(),
		deprecationTime.Day(),
		0, 0, 0, 0,
		losAngelesLocation,
	)

	// Convert deprecation time to client's timezone for comparison
	deprecationTimeClient := deprecationTimeLA.In(clientNow.Location())

	// Check if either the client time or LA time has passed the deprecation date
	// We use "either" to handle edge cases around timezone boundaries
	clientPassed := clientNow.After(deprecationTimeClient) || clientNow.Equal(deprecationTimeClient)
	laPassed := losAngelesNow.After(deprecationTimeLA) || losAngelesNow.Equal(deprecationTimeLA)

	return clientPassed || laPassed
}

// IsFrontDoorDeprecatedForCreation checks if Front Door creation is deprecated
func IsFrontDoorDeprecatedForCreation() bool {
	// Front Door was deprecated for new resource creation on April 1, 2025
	return isAPIFieldDeprecated("2025-04-01")
}

// IsFrontDoorFullyRetired checks if Front Door API is fully retired
func IsFrontDoorFullyRetired() bool {
	// Front Door API will be fully retired on March 31, 2027
	return isAPIFieldDeprecated("2027-03-31")
}
