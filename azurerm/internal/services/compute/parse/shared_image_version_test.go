package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SharedImageVersionId{}

func TestSharedImageVersionIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	galleryId := NewSharedImageGalleryId("group1", "gallery1")
	imageId := NewSharedImageId(galleryId, "image1")
	actual := NewSharedImageVersionId(imageId, "version1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/galleries/gallery1/images/image1/versions/version1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSharedImageVersionID(t *testing.T) {
	testData := []struct {
		Name   string
		Input  string
		Error  bool
		Expect *SharedImageVersionId
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/images/image1/versions/1.0.0",
			Error: true,
		},
		{
			Name:  "Missing image segment",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Compute/galleries/gallery1/versions/1.0.0",
			Error: true,
		},
		{
			Name:  "Image ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1/versions/1.0.0",
			Error: false,
			Expect: &SharedImageVersionId{
				ResourceGroup: "mygroup1",
				Gallery:       "gallery1",
				ImageName:     "image1",
				Version:       "1.0.0",
			},
		},
		{
			Name:  "Wrong Casing",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/galleries/gallery1/images/image1/Versions/1.0.0",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := SharedImageVersionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.ImageName != v.Expect.ImageName {
			t.Fatalf("Expected %q but got %q for Name", v.Expect.ImageName, actual.ImageName)
		}

		if actual.ResourceGroup != v.Expect.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expect.ResourceGroup, actual.ResourceGroup)
		}
	}
}
