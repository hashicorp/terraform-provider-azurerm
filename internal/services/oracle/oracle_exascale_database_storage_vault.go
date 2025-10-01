// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbstoragevaults"
)

func FlattenHighCapacityDatabaseStorage(input *exascaledbstoragevaults.ExascaleDbStorageDetails) []ExascaleDatabaseStorageDetailsModel {
	output := make([]ExascaleDatabaseStorageDetailsModel, 0)
	if input != nil {
		return append(output, ExascaleDatabaseStorageDetailsModel{
			AvailableSizeInGb: pointer.From(input.AvailableSizeInGbs),
			TotalSizeInGb:     pointer.From(input.TotalSizeInGbs),
		})
	}
	return output
}
