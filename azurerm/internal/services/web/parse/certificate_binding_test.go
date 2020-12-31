package parse

import (
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = CertificateBindingId{}

func TestCertificateBindingIDFormatter(t *testing.T) {
	hostnameBindingId := HostnameBindingId{
		SubscriptionId: "12345678-1234-9876-4563-123456789012",
		ResourceGroup:  "mygroup1",
		SiteName:       "site1",
		Name:           "binding1",
	}
	certificateId := CertificateId{
		SubscriptionId: "12345678-1234-9876-4563-123456789012",
		ResourceGroup:  "resGroup1",
		Name:           "certificate1",
	}
	actual := NewCertificateBindingId(hostnameBindingId, certificateId)
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/certificate1"
	if actual.ID() != expected {
		t.Fatalf("Expected %q, got %q", expected, actual)
	}
}

func TestCertificateBindingID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *CertificateBindingId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// Missing CertificateId Portion
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1|",
			Error: true,
		},

		{
			// Missing HostnameBindingId Portion
			Input: "|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/certificate1",
			Error: true,
		},

		{
			// Incorrect ordering
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/certificate1|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1",
			Error: true,
		},

		{
			// missing certificate name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/",
			Error: true,
		},

		{
			// invalid HostnameBindingId
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/certificate1",
			Error: true,
		},

		{
			// missing certificate name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1/hostNameBindings/binding1|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/certificates/certificate1",
			Expected: &CertificateBindingId{
				HostnameBindingId: HostnameBindingId{
					SubscriptionId: "12345678-1234-9876-4563-123456789012",
					ResourceGroup:  "mygroup1",
					SiteName:       "site1",
					Name:           "binding1",
				},
				CertificateId: CertificateId{
					SubscriptionId: "12345678-1234-9876-4563-123456789012",
					ResourceGroup:  "resGroup1",
					Name:           "certificate1",
				},
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := CertificateBindingID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value, but got an error: %s", err)
		}
		if v.Error {
			t.Fatalf("Expected an error, but did not get one")
		}

		if !reflect.DeepEqual(actual.HostnameBindingId, v.Expected.HostnameBindingId) {
			t.Fatalf("expected %+v, got %+v", v.Expected.HostnameBindingId, actual.HostnameBindingId)
		}

		if !reflect.DeepEqual(actual.CertificateId, v.Expected.CertificateId) {
			t.Fatalf("expected %+v, got %+v", v.Expected.CertificateId, actual.CertificateId)
		}
	}
}
