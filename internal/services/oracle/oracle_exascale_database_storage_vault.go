package oracle

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exascaledbstoragevaults"
)

type ExascaleDbStorageDetailsModel struct {
	AvailableSizeInGb int64 `tfschema:"available_size_in_gb"`
	TotalSizeInGb     int64 `tfschema:"total_size_in_gb"`
}

func flattenHighCapacityDatabaseStorage(input *exascaledbstoragevaults.ExascaleDbStorageDetails) []ExascaleDbStorageDetailsModel {
	output := make([]ExascaleDbStorageDetailsModel, 0)
	if input != nil {
		return append(output, ExascaleDbStorageDetailsModel{
			AvailableSizeInGb: pointer.From(input.AvailableSizeInGbs),
			TotalSizeInGb:     pointer.From(input.TotalSizeInGbs),
		})
	}
	return output
}
