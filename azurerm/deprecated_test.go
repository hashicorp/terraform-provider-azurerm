package azurerm

import (
	"regexp"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

// NOTE: these methods are deprecated, but provided to ease compatibility for open PR's

func testAccPreCheck(t *testing.T) {
	acceptance.PreCheck(t)
}

func testLocation() string {
	return acceptance.Location()
}

func testAltLocation() string {
	return acceptance.AltLocation()
}

func testAltLocation2() string {
	return acceptance.AltLocation2()
}

func testArmEnvironment() (*azure.Environment, error) {
	return acceptance.Environment()
}

func testGetAzureConfig(t *testing.T) *authentication.Config {
	return acceptance.GetAuthConfig(t)
}

func testRequiresImportError(resourceName string) *regexp.Regexp {
	return acceptance.RequiresImportError(resourceName)
}
