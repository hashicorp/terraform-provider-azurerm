// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestValidateAzureRmCosmosDbIndexingPolicy(t *testing.T) {
	cases := []struct {
		Name        string
		Value       *cosmosdb.IndexingPolicy
		ExpectError bool
	}{
		{
			Name:        "nil",
			Value:       nil,
			ExpectError: false,
		},
		{
			Name: "no included_path or excluded_path with Consistent indexing_mode",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeConsistent),
			},
			ExpectError: false,
		},
		{
			Name: "no included_path or excluded_path with None indexing_mode",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeNone),
			},
			ExpectError: false,
		},
		{
			Name: "included_path with /*",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeConsistent),
				IncludedPaths: &[]cosmosdb.IncludedPath{
					{
						Path: utils.String("/*"),
					},
					{
						Path: utils.String("/foo/?"),
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "excluded_path with /*",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeConsistent),
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: utils.String("/*"),
					},
					{
						Path: utils.String("/foo/?"),
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "included_path with /* and excluded_path",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeConsistent),
				IncludedPaths: &[]cosmosdb.IncludedPath{
					{
						Path: utils.String("/*"),
					},
					{
						Path: utils.String("/foo/?"),
					},
				},
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: utils.String("/testing/?"),
					},
					{
						Path: utils.String("/bar/?"),
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "included_path and excluded_path with /*",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeConsistent),
				IncludedPaths: &[]cosmosdb.IncludedPath{
					{
						Path: utils.String("/*"),
					},
					{
						Path: utils.String("/foo/?"),
					},
				},
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: utils.String("/*"),
					},
					{
						Path: utils.String("/testing/?"),
					},
					{
						Path: utils.String("/bar/?"),
					},
				},
			},
			ExpectError: true,
		},
		{
			Name: "missing /* from included_path",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeConsistent),
				IncludedPaths: &[]cosmosdb.IncludedPath{
					{
						Path: utils.String("/testing/?"),
					},
					{
						Path: utils.String("/foo/?"),
					},
				},
			},
			ExpectError: true,
		},
		{
			Name: "missing /* with included_path and excluded_path",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeConsistent),
				IncludedPaths: &[]cosmosdb.IncludedPath{
					{
						Path: utils.String("/foo/?"),
					},
					{
						Path: utils.String("/foo/?"),
					},
				},
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: utils.String("/bar/?"),
					},
					{
						Path: utils.String("/foo/?"),
					},
				},
			},
			ExpectError: true,
		},
		{
			Name: "indexing_mode None with included_path",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeNone),
				IncludedPaths: &[]cosmosdb.IncludedPath{
					{
						Path: utils.String("/*"),
					},
				},
			},
			ExpectError: true,
		},
		{
			Name: "indexing_mode None with excluded_path",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeNone),
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: utils.String("/*"),
					},
				},
			},
			ExpectError: true,
		},
		{
			Name: "indexing_mode None with included_path and excluded_path",
			Value: &cosmosdb.IndexingPolicy{
				IndexingMode: pointer.To(cosmosdb.IndexingModeNone),
				IncludedPaths: &[]cosmosdb.IncludedPath{
					{
						Path: utils.String("/*"),
					},
				},
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: utils.String("/testing/?"),
					},
				},
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		err := ValidateAzureRmCosmosDbIndexingPolicy(tc.Value)
		if tc.ExpectError && err == nil {
			t.Fatalf("Expected an error but didn't get one for %q", tc.Name)
		}

		if !tc.ExpectError && err != nil {
			t.Fatalf("Expected to get no errors for %q but got error: %+v", tc.Name, err)
		}
	}
}
