// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"slices"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
				KeyKind: pointer.ToEnum[cosmosdb.Kind]("NotAKeyKind"),
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

func TestFlattenCosmosDBAccountConnectionStrings(t *testing.T) {
	connectionString := func(description string, connectionType cosmosdb.Type, keyKind cosmosdb.Kind, value string) cosmosdb.DatabaseAccountConnectionString {
		return cosmosdb.DatabaseAccountConnectionString{
			ConnectionString: pointer.To(value),
			Description:      pointer.To(description),
			KeyKind:          pointer.To(keyKind),
			Type:             pointer.To(connectionType),
		}
	}

	sqlEntries := []cosmosdb.DatabaseAccountConnectionString{
		connectionString("Primary SQL Connection String", cosmosdb.TypeSql, cosmosdb.KindPrimary, "sql-primary"),
		connectionString("Secondary SQL Connection String", cosmosdb.TypeSql, cosmosdb.KindSecondary, "sql-secondary"),
		connectionString("Primary Read-Only SQL Connection String", cosmosdb.TypeSql, cosmosdb.KindPrimaryReadonly, "sql-primary-readonly"),
		connectionString("Secondary Read-Only SQL Connection String", cosmosdb.TypeSql, cosmosdb.KindSecondaryReadonly, "sql-secondary-readonly"),
	}

	primaryTableWithoutDescription := connectionString("", cosmosdb.TypeTable, cosmosdb.KindPrimary, "table-primary")
	primaryTableWithoutDescription.Description = nil

	tableEntries := []cosmosdb.DatabaseAccountConnectionString{
		primaryTableWithoutDescription,
		connectionString("Secondary Table Connection String", cosmosdb.TypeTable, cosmosdb.KindSecondary, "table-secondary"),
		connectionString("Primary Read-Only Table Connection String", cosmosdb.TypeTable, cosmosdb.KindPrimaryReadonly, "table-primary-readonly"),
		connectionString("Secondary Read-Only Table Connection String", cosmosdb.TypeTable, cosmosdb.KindSecondaryReadonly, "table-secondary-readonly"),
	}

	emptyTableAttributes := map[string]string{
		"primary_table_connection_string":            "",
		"secondary_table_connection_string":          "",
		"primary_readonly_table_connection_string":   "",
		"secondary_readonly_table_connection_string": "",
	}

	testData := []struct {
		name     string
		input    *[]cosmosdb.DatabaseAccountConnectionString
		expected map[string]string
	}{
		{
			name:  "table account",
			input: pointer.To(slices.Concat(sqlEntries, tableEntries)),
			expected: map[string]string{
				"primary_sql_connection_string":              "sql-primary",
				"secondary_sql_connection_string":            "sql-secondary",
				"primary_readonly_sql_connection_string":     "sql-primary-readonly",
				"secondary_readonly_sql_connection_string":   "sql-secondary-readonly",
				"primary_table_connection_string":            "table-primary",
				"secondary_table_connection_string":          "table-secondary",
				"primary_readonly_table_connection_string":   "table-primary-readonly",
				"secondary_readonly_table_connection_string": "table-secondary-readonly",
			},
		},
		{
			name:  "sql account clears the table attributes",
			input: pointer.To(sqlEntries),
			expected: map[string]string{
				"primary_sql_connection_string":              "sql-primary",
				"secondary_sql_connection_string":            "sql-secondary",
				"primary_readonly_sql_connection_string":     "sql-primary-readonly",
				"secondary_readonly_sql_connection_string":   "sql-secondary-readonly",
				"primary_table_connection_string":            "",
				"secondary_table_connection_string":          "",
				"primary_readonly_table_connection_string":   "",
				"secondary_readonly_table_connection_string": "",
			},
		},
		{
			name: "table entry without a connection string",
			input: pointer.To([]cosmosdb.DatabaseAccountConnectionString{
				{
					Type:    pointer.To(cosmosdb.TypeTable),
					KeyKind: pointer.To(cosmosdb.KindPrimary),
				},
			}),
			expected: emptyTableAttributes,
		},
		{
			name:     "no connection strings leaves sql and mongodb untouched",
			input:    nil,
			expected: emptyTableAttributes,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			actual := flattenCosmosDBAccountConnectionStrings(v.input)

			if len(actual) != len(v.expected) {
				t.Fatalf("expected %d attributes but got %d: %v", len(v.expected), len(actual), actual)
			}

			for key, expected := range v.expected {
				value, ok := actual[key]
				if !ok {
					t.Fatalf("expected attribute %q to be returned", key)
				}
				if value != expected {
					t.Fatalf("expected %q to be %q but got %q", key, expected, value)
				}
			}
		})
	}
}

func TestFlattenCosmosDBAccountConnectionStringsMatchSchema(t *testing.T) {
	input := make([]cosmosdb.DatabaseAccountConnectionString, 0)
	for description := range connStringPropertyMap {
		input = append(input, cosmosdb.DatabaseAccountConnectionString{
			Description: pointer.To(description),
		})
	}
	for _, keyKind := range cosmosdb.PossibleValuesForKind() {
		input = append(input, cosmosdb.DatabaseAccountConnectionString{
			KeyKind: pointer.ToEnum[cosmosdb.Kind](keyKind),
			Type:    pointer.To(cosmosdb.TypeTable),
		})
	}

	attributes := flattenCosmosDBAccountConnectionStrings(&input)
	if expected := len(connStringPropertyMap) + len(cosmosdb.PossibleValuesForKind()); len(attributes) != expected {
		t.Fatalf("expected %d attributes but got %d: %v", expected, len(attributes), attributes)
	}

	schemas := map[string]map[string]*pluginsdk.Schema{
		"resource":    resourceCosmosDbAccount().Schema,
		"data source": dataSourceCosmosDbAccount().Schema,
	}

	for name, schema := range schemas {
		for attribute := range attributes {
			v, ok := schema[attribute]
			if !ok {
				t.Fatalf("the %s schema has no %q attribute", name, attribute)
			}
			if v.Type != pluginsdk.TypeString || !v.Computed || !v.Sensitive {
				t.Fatalf("the %s attribute %q should be a computed, sensitive string", name, attribute)
			}
		}
	}
}
