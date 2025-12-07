// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"time"
)

var (
	CreateDeprecationMessage  string = "creation of new CDN resources is no longer permitted following its deprecation on October 1, 2025. However, modifications to existing CDN resources remain supported until the API reaches full retirement on September 30, 2027"
	VerizonDeprecationMessage string = "creation of CDN resources are no longer permitted with the `StandardVerizon` and `PremiumVerizon` sku's following its deprecation on October 1, 2025"
	AkamaiDeprecationMessage  string = "creation of CDN resources are no longer permitted with the `StandardAkamai` sku following its deprecation on October 31, 2023"
	FullyRetiredMessage       string = "as of September 30, 2027, CDN resources have reached full retirement and are no longer supported by Azure"
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

// IsCdnDeprecatedForCreation checks if CDN creation is deprecated
func IsCdnDeprecatedForCreation() bool {
	// Cdn will be deprecated for new resource creation on October 1, 2025
	return isAPIFieldDeprecated("2025-10-01")
}

// IsCdnFullyRetired checks if CDB API is fully retired
func IsCdnFullyRetired() bool {
	// Cdn API will be fully retired on September 30, 2027
	return isAPIFieldDeprecated("2027-09-30")
}

// IsCdnStandardAkamaiiDeprecated checks if CDN sku StandardAkamai is deprecated for new resource creation
func IsCdnStandardAkamaiDeprecatedForCreation() bool {
	// Cdn sku StandardAkamai was deprecated for new resource creation on October 31, 2023
	return isAPIFieldDeprecated("2023-10-31")
}

// IsCdnStandardVerizonPremiumVerizonDeprecatedForCreation checks if CDN sku StandardVerizon and PremiumVerizon is deprecated for new resource creation
func IsCdnStandardVerizonPremiumVerizonDeprecatedForCreation() bool {
	// Cdn sku StandardVerizon and PremiumVerizon was deprecated for new resource creation on January 15, 2025
	return isAPIFieldDeprecated("2025-01-15")
}
