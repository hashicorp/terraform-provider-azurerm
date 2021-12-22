package compute

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestSortVersions(t *testing.T) {
	testData := []struct {
		input    []compute.GalleryImageVersion
		expected []compute.GalleryImageVersion
	}{
		{
			input: []compute.GalleryImageVersion{
				{Name: utils.String("1.0.1")},
				{Name: utils.String("1.2.15.0")},
				{Name: utils.String("1.0.8")},
				{Name: utils.String("1.0.9")},
				{Name: utils.String("1.0.1.1")},
				{Name: utils.String("1.0.10")},
			},
			expected: []compute.GalleryImageVersion{
				{Name: utils.String("1.0.1")},
				{Name: utils.String("1.0.1.1")},
				{Name: utils.String("1.0.8")},
				{Name: utils.String("1.0.9")},
				{Name: utils.String("1.0.10")},
				{Name: utils.String("1.2.15.0")},
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing sortSharedImageVersions..")

		actual := sortSharedImageVersions(v.input)
		if eq := reflect.DeepEqual(v.expected, actual); !eq {
			t.Fatalf("Expected %v but got %v", v.expected, actual)
		}
	}
}
