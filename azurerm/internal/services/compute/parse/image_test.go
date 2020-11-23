package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ImageId{}

func TestImageIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewImageId(subscriptionId, "group1", "image1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/images/image1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestImageID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *ImageId
	}{
		{
			Name:  "Empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "No Resource Groups Segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			Name:  "No Resource Groups Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},
		{
			Name:  "Resource Group ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Error: true,
		},
		{
			Name:  "Missing Image Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/images/",
			Error: true,
		},
		{
			Name:  "Image ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/images/image1",
			Error: false,
			Expect: &ImageId{
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
				ResourceGroup:  "resGroup1",
				Name:           "image1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/Images/image1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ImageID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
