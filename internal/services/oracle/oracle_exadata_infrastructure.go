// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudexadatainfrastructures"
)

type ExascaleConfigDetails struct {
	TotalStorageInGb     int64 `tfschema:"total_storage_in_gb"`
	AvailableStorageInGb int64 `tfschema:"available_storage_in_gb"`
}

func FlattenExascaleConfig(input *cloudexadatainfrastructures.ExascaleConfigDetails) []ExascaleConfigDetails {
	output := make([]ExascaleConfigDetails, 0)
	if input != nil {
		output = append(output, ExascaleConfigDetails{
			TotalStorageInGb:     input.TotalStorageInGbs,
			AvailableStorageInGb: pointer.From(input.AvailableStorageInGbs),
		})
	}
	return output
}
