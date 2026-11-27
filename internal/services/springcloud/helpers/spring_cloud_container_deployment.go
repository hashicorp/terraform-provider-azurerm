// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func ExpandSpringCloudDeploymentApms(input []any) *[]appplatform.ApmReference {
	if len(input) == 0 {
		return nil
	}
	result := make([]appplatform.ApmReference, 0)
	for _, v := range input {
		result = append(result, appplatform.ApmReference{
			ResourceID: pointer.To(v.(string)),
		})
	}
	return pointer.To(result)
}

func FlattenSpringCloudDeploymentApms(input *[]appplatform.ApmReference) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}
	result := make([]any, 0)
	for _, v := range *input {
		id, err := appplatform_rm.ParseApmIDInsensitively(*v.ResourceID)
		if err != nil {
			return nil, err
		}
		result = append(result, id.ID())
	}
	return result, nil
}
