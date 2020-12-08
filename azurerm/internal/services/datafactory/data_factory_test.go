package datafactory

import "testing"

func TestAzureRmDataFactoryLinkedServiceConnectionStringDiff(t *testing.T) {
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

func TestAzureRmDataFactoryDeserializePipelineActivities(t *testing.T) {
	cases := []struct {
		Json                string
		ExpectActivityCount int
		ExpectErr           bool
	}{
		{
			Json:                "{}",
			ExpectActivityCount: 0,
			ExpectErr:           true,
		},
		{
			Json: `[
				{
				  "type": "ForEach",
				  "typeProperties": {
					"isSequential": true,
					"items": {
					  "value": "@pipeline().parameters.OutputBlobNameList",
					  "type": "Expression"
					},
					"activities": [
					  {
						"type": "Copy",
						"typeProperties": {
						  "source": {
							"type": "BlobSource"
						  },
						  "sink": {
							"type": "BlobSink"
						  },
						  "dataIntegrationUnits": 32
						},
						"inputs": [
						  {
							"referenceName": "exampleDataset",
							"parameters": {
							  "MyFolderPath": "examplecontainer",
							  "MyFileName": "examplecontainer.csv"
							},
							"type": "DatasetReference"
						  }
						],
						"outputs": [
						  {
							"referenceName": "exampleDataset",
							"parameters": {
							  "MyFolderPath": "examplecontainer",
							  "MyFileName": {
								"value": "@item()",
								"type": "Expression"
							  }
							},
							"type": "DatasetReference"
						  }
						],
						"name": "ExampleCopyActivity"
					  }
					]
				  },
				  "name": "ExampleForeachActivity"
				}
			  ]`,
			ExpectActivityCount: 1,
			ExpectErr:           false,
		},
	}

	for _, tc := range cases {
		items, err := deserializeDataFactoryPipelineActivities(tc.Json)
		if err != nil {
			if tc.ExpectErr {
				t.Log("Expected error and got error")
				return
			}

			t.Fatal(err)
		}

		if items == nil && !tc.ExpectErr {
			t.Fatal("Expected items got nil")
		}

		if len(*items) != tc.ExpectActivityCount {
			t.Fatal("Failed to deserialise pipeline")
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
