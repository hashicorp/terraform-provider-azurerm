// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
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
						Path: pointer.To("/*"),
					},
					{
						Path: pointer.To("/foo/?"),
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
						Path: pointer.To("/*"),
					},
					{
						Path: pointer.To("/foo/?"),
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
						Path: pointer.To("/*"),
					},
					{
						Path: pointer.To("/foo/?"),
					},
				},
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: pointer.To("/testing/?"),
					},
					{
						Path: pointer.To("/bar/?"),
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
						Path: pointer.To("/*"),
					},
					{
						Path: pointer.To("/foo/?"),
					},
				},
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: pointer.To("/*"),
					},
					{
						Path: pointer.To("/testing/?"),
					},
					{
						Path: pointer.To("/bar/?"),
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
						Path: pointer.To("/testing/?"),
					},
					{
						Path: pointer.To("/foo/?"),
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
						Path: pointer.To("/foo/?"),
					},
					{
						Path: pointer.To("/foo/?"),
					},
				},
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: pointer.To("/bar/?"),
					},
					{
						Path: pointer.To("/foo/?"),
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
						Path: pointer.To("/*"),
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
						Path: pointer.To("/*"),
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
						Path: pointer.To("/*"),
					},
				},
				ExcludedPaths: &[]cosmosdb.ExcludedPath{
					{
						Path: pointer.To("/testing/?"),
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
