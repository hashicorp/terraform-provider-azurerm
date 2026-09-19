// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
)

func TestTableConnectionStringAttribute(t *testing.T) {
	testData := []struct {
		name     string
		input    cosmosdb.DatabaseAccountConnectionString
		expected string
	}{
		{
			name: "primary table",
			input: cosmosdb.DatabaseAccountConnectionString{
				Type:    pointer.To(cosmosdb.TypeTable),
				KeyKind: pointer.To(cosmosdb.KindPrimary),
			},
			expected: "primary_table_connection_string",
		},
		{
			name: "secondary table",
			input: cosmosdb.DatabaseAccountConnectionString{
				Type:    pointer.To(cosmosdb.TypeTable),
				KeyKind: pointer.To(cosmosdb.KindSecondary),
			},
			expected: "secondary_table_connection_string",
		},
		{
			name: "primary read-only table",
			input: cosmosdb.DatabaseAccountConnectionString{
				Type:    pointer.To(cosmosdb.TypeTable),
				KeyKind: pointer.To(cosmosdb.KindPrimaryReadonly),
			},
			expected: "primary_readonly_table_connection_string",
		},
		{
			name: "secondary read-only table",
			input: cosmosdb.DatabaseAccountConnectionString{
				Type:    pointer.To(cosmosdb.TypeTable),
				KeyKind: pointer.To(cosmosdb.KindSecondaryReadonly),
			},
			expected: "secondary_readonly_table_connection_string",
		},
		{
			name: "sql is left to connStringPropertyMap",
			input: cosmosdb.DatabaseAccountConnectionString{
				Type:    pointer.To(cosmosdb.TypeSql),
				KeyKind: pointer.To(cosmosdb.KindPrimary),
			},
			expected: "",
		},
		{
			name: "mongodb is left to connStringPropertyMap",
			input: cosmosdb.DatabaseAccountConnectionString{
				Type:    pointer.To(cosmosdb.TypeMongoDB),
				KeyKind: pointer.To(cosmosdb.KindPrimary),
			},
			expected: "",
		},
		{
			name: "table with an unknown key kind",
			input: cosmosdb.DatabaseAccountConnectionString{
				Type:    pointer.To(cosmosdb.TypeTable),
				KeyKind: pointer.To(cosmosdb.Kind("NotAKeyKind")),
			},
			expected: "",
		},
		{
			name: "table with no key kind",
			input: cosmosdb.DatabaseAccountConnectionString{
				Type: pointer.To(cosmosdb.TypeTable),
			},
			expected: "",
		},
		{
			name:     "no type and no key kind",
			input:    cosmosdb.DatabaseAccountConnectionString{},
			expected: "",
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			if actual := tableConnectionStringAttribute(v.input); actual != v.expected {
				t.Fatalf("expected %q but got %q", v.expected, actual)
			}
		})
	}
}
