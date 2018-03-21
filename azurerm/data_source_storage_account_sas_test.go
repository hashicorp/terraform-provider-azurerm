package azurerm

import (
	"fmt"
	"testing"
)

func TestDataSourceArmStorageAccountSasRead(t *testing.T) {

}

func testAccDataSourceAzureRMStorageAccountSas_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestsa-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name = "acctestsads%s"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location = "${azurerm_resource_group.test.location}"
  account_tier = "Standard"
  account_replication_type = "LRS"

  tags {
    environment = "production"
  }
}

data "azurerm_storage_account_sas" "test" {
  connection_string = "${azurerm_storage_account.test.primary_connection_string}"
  https_only        = true
  resource_types {
    service   = true
    container = false
    object    = false
  }
  services {
    blob  = true
    queue = false
    table = false
	file  = false
  }
  start   = ""
  expirty = ""
  permissions {
	read    = true
	write   = true
	delete  = false
    list    = false
    add     = true
	create  = true
	update  = false
	process = false
  }
}
`, rInt, location, rString)
}

func TestBuildResourceTypesString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{"service": true}, "s"},
		{map[string]interface{}{"container": true}, "c"},
		{map[string]interface{}{"object": true}, "o"},
		{map[string]interface{}{"service": true, "container": true, "object": true}, "sco"},
	}

	for _, test := range testCases {
		result := buildResourceTypesString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}

func TestBuildServicesString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{"blob": true}, "b"},
		{map[string]interface{}{"queue": true}, "q"},
		{map[string]interface{}{"table": true}, "t"},
		{map[string]interface{}{"file": true}, "f"},
		{map[string]interface{}{"blob": true, "queue": true, "table": true, "file": true}, "bqtf"},
	}

	for _, test := range testCases {
		result := buildServicesString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}

}

func TestBuildPermissionsString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{"read": true}, "r"},
		{map[string]interface{}{"write": true}, "w"},
		{map[string]interface{}{"delete": true}, "d"},
		{map[string]interface{}{"list": true}, "l"},
		{map[string]interface{}{"add": true}, "a"},
		{map[string]interface{}{"create": true}, "c"},
		{map[string]interface{}{"update": true}, "u"},
		{map[string]interface{}{"process": true}, "p"},
		{map[string]interface{}{"read": true, "write": true, "add": true, "create": true}, "rwac"},
	}

	for _, test := range testCases {
		result := buildPermissionsString(test.input)
		if test.expected != result {
			t.Fatalf("Failed to build resource type string: expected: %s, result: %s", test.expected, result)
		}
	}
}

func TestValidateArmStorageAccountSasResourceTypes(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"s", false},
		{"c", false},
		{"o", false},
		{"sc", false},
		{"cs", false},
		{"os", false},
		{"sco", false},
		{"cso", false},
		{"osc", false},
		{"scos", true},
		{"csoc", true},
		{"oscs", true},
		{"S", true},
		{"C", true},
		{"O", true},
	}

	for _, test := range testCases {
		_, es := validateArmStorageAccountSasResourceTypes(test.input, "<unused>")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating resource_types %q to fail", test.input)
		}
	}
}

// This connection string was for a real storage account which has been deleted
// so its safe to include here for reference to understand the format.
// DefaultEndpointsProtocol=https;AccountName=azurermtestsa0;AccountKey=T0ZQouXBDpWud/PlTRHIJH2+VUK8D+fnedEynb9Mx638IYnsMUe4mv1fFjC7t0NayTfFAQJzPZuV1WHFKOzGdg==;EndpointSuffix=core.windows.net

func TestComputeAzureStorageAccountSas(t *testing.T) {
	testCases := []struct {
		accountName    string
		accountKey     string
		permissions    string
		services       string
		resourceTypes  string
		start          string
		expiry         string
		signedProtocol string
		signedIp       string
		signedVersion  string
		knownSasToken  string
	}{
		{
			"azurermtestsa0",
			"T0ZQouXBDpWud/PlTRHIJH2+VUK8D+fnedEynb9Mx638IYnsMUe4mv1fFjC7t0NayTfFAQJzPZuV1WHFKOzGdg==",
			"rwac",
			"b",
			"c",
			"2018-03-20T04:00:00Z",
			"2020-03-20T04:00:00Z",
			"https",
			"",
			"2017-07-29",
			"?sv=2017-07-29&ss=b&srt=c&sp=rwac&se=2020-03-20T04:00:00Z&st=2018-03-20T04:00:00Z&spr=https&sig=SQigK%2FnFA4pv0F0oMLqr6DxUWV4vtFqWi6q3Mf7o9nY%3D",
		},
		{
			"azurermtestsa0",
			"2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==",
			"rwdlac",
			"b",
			"sco",
			"2018-03-20T04:00:00Z",
			"2018-03-28T05:04:25Z",
			"https,http",
			"",
			"2017-07-29",
			"?sv=2017-07-29&ss=b&srt=sco&sp=rwdlac&se=2018-03-28T05:04:25Z&st=2018-03-20T04:00:00Z&spr=https,http&sig=OLNwL%2B7gxeDQQaUyNdXcDPK2aCbCMgEkJNjha9te448%3D",
		},
	}

	for _, test := range testCases {
		computedToken, err := computeAzureStorageAccountSas(test.accountName,
			test.accountKey,
			test.permissions,
			test.services,
			test.resourceTypes,
			test.start,
			test.expiry,
			test.signedProtocol,
			test.signedIp,
			test.signedVersion)

		if err != nil {
			t.Fatalf("Test Failed: Error computing storage account Sas: %q", err)
		}

		if computedToken != test.knownSasToken {
			t.Fatalf("Test failed: Expected Azure SAS %s but was %s", test.knownSasToken, computedToken)
		}
	}
}
