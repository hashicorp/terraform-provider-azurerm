package parse

import (
	"testing"
)

func TestBatchCertificateId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *CertificateId
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
			Name:     "Missing Accounts Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/",
			Expected: nil,
		},
		{
			Name:     "Batch Account ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/acctName1",
			Expected: nil,
		},
		{
			Name:     "Missing Certificates Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/acctName1/certificates/",
			Expected: nil,
		},
		{
			Name:  "Batch Certificate ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/acctName1/certificates/Certificate1",
			Expected: &CertificateId{
				ResourceGroup:    "resGroup1",
				BatchAccountName: "acctName1",
				CertificateName:  "Certificate1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/acctName1/Certificates/Certificate1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := CertificateID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.CertificateName != v.Expected.CertificateName {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.CertificateName, actual.CertificateName)
		}
		if actual.BatchAccountName != v.Expected.BatchAccountName {
			t.Fatalf("Expected %q but got %q for BatchAccountName", v.Expected.BatchAccountName, actual.BatchAccountName)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
