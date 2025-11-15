// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exascaledbstoragevaults"
)

type ExascaleDatabaseStorageDetailsModel struct {
	AvailableSizeInGb int64 `tfschema:"available_size_in_gb"`
	TotalSizeInGb     int64 `tfschema:"total_size_in_gb"`
}

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

func FlattenAttachedShapeAttribute(input *[]exascaledbstoragevaults.ShapeAttribute) []string {
	output := make([]string, 0)
	if input != nil {
		for _, value := range *input {
			output = append(output, string(value))
		}
	}
	return output
}
