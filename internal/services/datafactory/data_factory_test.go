// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import "testing"

func TestDataFactoryLinkedServiceConnectionStringDiff(t *testing.T) {
	cases := []struct {
		Old    string
		New    string
		NoDiff bool
	}{
		{
			Old:    "",
			New:    "",
			NoDiff: true,
		},
		{
			Old:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test",
			New:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test;Password=test",
			NoDiff: true,
		},
		{
			Old:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test",
			New:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test",
			NoDiff: true,
		},
		{
			Old:    "Integrated Security=False;Data Source=test2;Initial Catalog=test;User ID=test",
			New:    "Integrated Security=False;Data Source=test;Initial Catalog=test;User ID=test;Password=test",
			NoDiff: false,
		},
	}

	for _, tc := range cases {
		noDiff := azureRmDataFactoryLinkedServiceConnectionStringDiff("", tc.Old, tc.New, nil)

		if noDiff != tc.NoDiff {
			t.Fatalf("Expected azureRmDataFactoryLinkedServiceConnectionStringDiff to be '%t' for '%s' '%s' - got '%t'", tc.NoDiff, tc.Old, tc.New, noDiff)
		}
	}
}

func TestNormalizeJSON(t *testing.T) {
	cases := []struct {
		Old      string
		New      string
		Suppress bool
	}{
		{
			Old: `[
				{
					"name": "Append variable1",
					"type": "AppendVariable",
					"dependsOn": [],
					"userProperties": [],
					"typeProperties": {
						"variableName": "bob",
						"value": "something"
					}
				}
			]`,
			New: `[
				{
					"name": "Append variable1",
					"type": "AppendVariable",
					"dependsOn": [],
					"userProperties": [],
					"typeProperties": {
						"value": "something",
						"variableName": "bob"
					}
				}
			]`,
			Suppress: true,
		},
		{
			Old: `[
				{
					"name": "Append variable1",
					"type": "AppendVariable",
					"dependsOn": [],
					"userProperties": [],
					"typeProperties": {
						"variableName": "bobdifferent",
						"value": "something"
					}
				}
			]`,
			New: `[
				{
					"name": "Append variable1",
					"type": "AppendVariable",
					"dependsOn": [],
					"userProperties": [],
					"typeProperties": {
						"value": "something",
						"variableName": "bob"
					}
				}
			]`,
			Suppress: false,
		},
		{
			Old:      `{ "notbob": "notbill" }`,
			New:      `{ "bob": "bill" }`,
			Suppress: false,
		},
	}

	for _, tc := range cases {
		suppress := suppressJsonOrderingDifference("test", tc.Old, tc.New, nil)

		if suppress != tc.Suppress {
			t.Fatalf("Expected JsonOrderingDifference to be '%t' for '%s' '%s' - got '%t'", tc.Suppress, tc.Old, tc.New, suppress)
		}
	}
}

func TestExpandCompressionType(t *testing.T) {
	cases := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          "Gzip",
			expectedOutput: "gzip",
		},
		{
			input:          "gzip",
			expectedOutput: "gzip",
		},
		{
			input:          "TarGZip",
			expectedOutput: "TarGZip",
		},
	}

	for _, tc := range cases {
		output := expandCompressionType(tc.input)

		if output != tc.expectedOutput {
			t.Fatalf("Expected expandCompressionType to be '%s' for '%s' - got '%s'", tc.expectedOutput, tc.input, output)
		}
	}
}
