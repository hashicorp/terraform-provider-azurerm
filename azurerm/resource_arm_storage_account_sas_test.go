package azurerm

import (
	"testing"
)

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
// DefaultEndpointsProtocol=https;AccountName=azurermtestsa0;AccountKey=2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==;EndpointSuffix=core.windows.net

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
			"2vJrjEyL4re2nxCEg590wJUUC7PiqqrDHjAN5RU304FNUQieiEwS2bfp83O0v28iSfWjvYhkGmjYQAdd9x+6nw==",
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
	}

	for _, test := range testCases {
		computedToken := computeAzureStorageAccountSas(test.accountName,
			test.accountKey,
			test.permissions,
			test.services,
			test.resourceTypes,
			test.start,
			test.expiry,
			test.signedProtocol,
			test.signedIp,
			test.signedVersion)

		if computedToken != test.knownSasToken {
			t.Fatalf("Test failed: Expected Azure SAS %s but was %s", test.knownSasToken, computedToken)
		}
	}
}
