package videoanalyzer

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = EdgeModuleId{}

func TestEdgeModuleIDFormatter(t *testing.T) {
	actual := NewEdgeModuleID("{subscriptionId}", "{resourceGroupName}", "{accountName}", "{edgeModuleName}").ID()
	expected := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/{edgeModuleName}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseEdgeModuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *EdgeModuleId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/{subscriptionId}/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/{subscriptionId}/resourceGroups/",
			Error: true,
		},

		{
			// missing VideoAnalyzerName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/",
			Error: true,
		},

		{
			// missing value for VideoAnalyzerName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/{edgeModuleName}",
			Expected: &EdgeModuleId{
				SubscriptionId:    "{subscriptionId}",
				ResourceGroup:     "{resourceGroupName}",
				VideoAnalyzerName: "{accountName}",
				Name:              "{edgeModuleName}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/RESOURCEGROUPS/{RESOURCEGROUPNAME}/PROVIDERS/MICROSOFT.MEDIA/VIDEOANALYZERS/{ACCOUNTNAME}/EDGEMODULES/{EDGEMODULENAME}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseEdgeModuleID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.VideoAnalyzerName != v.Expected.VideoAnalyzerName {
			t.Fatalf("Expected %q but got %q for VideoAnalyzerName", v.Expected.VideoAnalyzerName, actual.VideoAnalyzerName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestParseEdgeModuleIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *EdgeModuleId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/{subscriptionId}/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/{subscriptionId}/resourceGroups/",
			Error: true,
		},

		{
			// missing VideoAnalyzerName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/",
			Error: true,
		},

		{
			// missing value for VideoAnalyzerName
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoAnalyzers/{accountName}/edgeModules/{edgeModuleName}",
			Expected: &EdgeModuleId{
				SubscriptionId:    "{subscriptionId}",
				ResourceGroup:     "{resourceGroupName}",
				VideoAnalyzerName: "{accountName}",
				Name:              "{edgeModuleName}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/videoanalyzers/{accountName}/edgemodules/{edgeModuleName}",
			Expected: &EdgeModuleId{
				SubscriptionId:    "{subscriptionId}",
				ResourceGroup:     "{resourceGroupName}",
				VideoAnalyzerName: "{accountName}",
				Name:              "{edgeModuleName}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/VIDEOANALYZERS/{accountName}/EDGEMODULES/{edgeModuleName}",
			Expected: &EdgeModuleId{
				SubscriptionId:    "{subscriptionId}",
				ResourceGroup:     "{resourceGroupName}",
				VideoAnalyzerName: "{accountName}",
				Name:              "{edgeModuleName}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/ViDeOaNaLyZeRs/{accountName}/EdGeMoDuLeS/{edgeModuleName}",
			Expected: &EdgeModuleId{
				SubscriptionId:    "{subscriptionId}",
				ResourceGroup:     "{resourceGroupName}",
				VideoAnalyzerName: "{accountName}",
				Name:              "{edgeModuleName}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseEdgeModuleIDInsensitively(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.VideoAnalyzerName != v.Expected.VideoAnalyzerName {
			t.Fatalf("Expected %q but got %q for VideoAnalyzerName", v.Expected.VideoAnalyzerName, actual.VideoAnalyzerName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
