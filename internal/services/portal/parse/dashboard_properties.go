// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/dashboards"
)

// rawDashboardPartMetadata is an adapter that implements dashboards.DashboardPartMetadata
// while preserving the original JSON structure. The Azure SDK's RawDashboardPartMetadataImpl
// wraps unrecognised part types in a {Type, Values} envelope on serialisation, which corrupts
// the round-trip for metadata that users paste from the Azure Portal. This adapter avoids the
// envelope by holding the raw map directly.
type rawDashboardPartMetadata map[string]interface{}

var _ dashboards.DashboardPartMetadata = rawDashboardPartMetadata{}

func (r rawDashboardPartMetadata) DashboardPartMetadata() dashboards.BaseDashboardPartMetadataImpl {
	metadataType := dashboards.DashboardPartMetadataType("")
	if v, ok := r["type"]; ok {
		metadataType = dashboards.DashboardPartMetadataType(fmt.Sprintf("%v", v))
	}

	return dashboards.BaseDashboardPartMetadataImpl{Type: metadataType}
}

// NormalizeDashboardPartMetadata replaces any RawDashboardPartMetadataImpl instances
// with rawDashboardPartMetadata so the metadata serialises without the {Type, Values}
// envelope that the Azure SDK adds for unrecognised discriminator types.
func NormalizeDashboardPartMetadata(props *dashboards.DashboardPropertiesWithProvisioningState) {
	if props == nil || props.Lenses == nil {
		return
	}

	for lensIndex := range *props.Lenses {
		for partIndex := range (*props.Lenses)[lensIndex].Parts {
			part := &(*props.Lenses)[lensIndex].Parts[partIndex]
			if part.Metadata == nil {
				continue
			}

			rawMetadata, ok := part.Metadata.(dashboards.RawDashboardPartMetadataImpl)
			if !ok || rawMetadata.Values == nil {
				continue
			}

			part.Metadata = rawDashboardPartMetadata(rawMetadata.Values)
		}
	}
}
