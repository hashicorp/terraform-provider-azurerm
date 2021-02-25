package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ImageBuilderTemplateId{}

func TestImageBuilderTemplateIDFormatter(t *testing.T) {
	actual := NewImageBuilderTemplateID("12345678-1234-9876-4563-123456789012", "resGroup1", "imageTemplate1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.VirtualMachineImages/imageTemplates/imageTemplate1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestImageBuilderTemplateID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ImageBuilderTemplateId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing ImageTemplateName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.VirtualMachineImages/",
			Error: true,
		},

		{
			// missing value for ImageTemplateName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.VirtualMachineImages/imageTemplates/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.VirtualMachineImages/imageTemplates/imageTemplate1",
			Expected: &ImageBuilderTemplateId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resGroup1",
				ImageTemplateName: "imageTemplate1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.VIRTUALMACHINEIMAGES/IMAGETEMPLATES/IMAGETEMPLATE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ImageBuilderTemplateID(v.Input)
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
		if actual.ImageTemplateName != v.Expected.ImageTemplateName {
			t.Fatalf("Expected %q but got %q for ImageTemplateName", v.Expected.ImageTemplateName, actual.ImageTemplateName)
		}
	}
}
