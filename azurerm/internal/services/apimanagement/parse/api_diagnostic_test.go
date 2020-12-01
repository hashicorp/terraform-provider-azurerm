package parse

import "testing"

func TestApiManagementApiDiagnosticID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *ApiDiagnosticId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Service Name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/",
			Expected: nil,
		},
		{
			Name:     "Missing APIs",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1",
			Expected: nil,
		},
		{
			Name:     "Missing APIs Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis",
			Expected: nil,
		},
		{
			Name:     "Missing Diagnostics",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1",
			Expected: nil,
		},
		{
			Name:     "Missing Diagnostics Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1/diagnostics",
			Expected: nil,
		},
		{
			Name:  "Diagnostic ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/apis/api1/diagnostics/applicationinsights",
			Expected: &ApiDiagnosticId{
				DiagnosticName: "applicationinsights",
				ApiName:        "api1",
				ServiceName:    "service1",
				ResourceGroup:  "resGroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/APIs/api1/diagnostics/applicationinsights",
			Expected: nil,
		},
		{
			Name:  "From ACC test",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG-200904094049117016/providers/Microsoft.ApiManagement/service/acctestAM-200904094049117016/apis/acctestAMA-200904094049117016/diagnostics/applicationinsights",
			Expected: &ApiDiagnosticId{
				DiagnosticName: "applicationinsights",
				ApiName:        "acctestAMA-200904094049117016",
				ServiceName:    "acctestAM-200904094049117016",
				ResourceGroup:  "acctestRG-200904094049117016",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ApiDiagnosticID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.DiagnosticName != v.Expected.DiagnosticName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.DiagnosticName, actual.DiagnosticName)
		}

		if actual.ApiName != v.Expected.ApiName {
			t.Fatalf("Expected %q but got %q for API Name", v.Expected.ApiName, actual.ApiName)
		}

		if actual.ServiceName != v.Expected.ServiceName {
			t.Fatalf("Expected %q but got %q for Service Name", v.Expected.ServiceName, actual.ServiceName)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
