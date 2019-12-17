package acceptance

import (
	"fmt"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

type TestData struct {
	// Locations is a set of Azure Regions which should be used for this Test
	Locations Regions

	// RandomString is a random integer which is unique to this test case
	RandomInteger int

	// RandomString is a random 5 character string is unique to this test case
	RandomString string

	// ResourceName is the fully qualified resource name, comprising of the
	// resource type and then the resource label
	// e.g. `azurerm_resource_group.test`
	ResourceName string

	// Environment is a struct containing Details about the Azure Environment
	// that we're running against
	Environment azure.Environment

	// EnvironmentName is the name of the Azure Environment where we're running
	EnvironmentName string

	// resourceType is the Terraform Resource Type - `azurerm_resource_group`
	resourceType string

	// resourceLabel is the local used for the resource - generally "test""
	resourceLabel string
}

// BuildTestData generates some test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	env, err := Environment()
	if err != nil {
		t.Fatalf("Error retrieving Environment: %+v", err)
	}

	testData := TestData{
		RandomInteger:   tf.AccRandTimeInt(),
		RandomString:    acctest.RandString(5),
		ResourceName:    fmt.Sprintf("%s.%s", resourceType, resourceLabel),
		Environment:     *env,
		EnvironmentName: EnvironmentName(),

		resourceType:  resourceType,
		resourceLabel: resourceLabel,
	}

	if features.UseDynamicTestLocations() {
		testData.Locations = availableLocations()
	} else {
		testData.Locations = Regions{
			Primary:   Location(),
			Secondary: AltLocation(),
			Ternary:   AltLocation2(),
		}
	}

	return testData
}
