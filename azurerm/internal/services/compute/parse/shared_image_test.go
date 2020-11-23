package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SharedImageId{}

func TestSharedImageIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	galleryId := NewSharedImageGalleryId(subscriptionId, "group1", "gallery1")
	actual := NewSharedImageId(galleryId, "image1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/galleries/gallery1/images/image1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSharedImageID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *SharedImageId
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
			Name:  "Missing galleries segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/images/image1",
			Error: true,
		},
		{
			Name:  "Missing image Value",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/galleries/gallery1/images",
			Error: true,
		},
		{
			Name:  "Image ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1",
			Error: false,
			Expect: &SharedImageId{
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
				ResourceGroup:  "mygroup1",
				Gallery:        "gallery1",
				Name:           "image1",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1/Images/image1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SharedImageID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expect.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.Name, actual.Name)
		}

		if actual.Gallery != v.Expect.Gallery {
			t.Fatalf("Expected %q but got %q for Gallery", v.Expect.Gallery, actual.Gallery)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}

		if actual.SubscriptionId != v.Expect.SubscriptionId {
			t.Fatalf("Expected %q but got %q for Subscription Id", v.Expect.SubscriptionId, actual.SubscriptionId)
		}
	}
}
