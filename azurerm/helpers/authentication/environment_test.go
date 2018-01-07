package authentication

import (
	"testing"
)

func TestAzureEnvironmentNames(t *testing.T) {
	testData := map[string]string{
		"":                       "public",
		"AzureChinaCloud":        "china",
		"AzureCloud":             "public",
		"AzureGermanCloud":       "german",
		"AZUREUSGOVERNMENTCLOUD": "usgovernment",
		"AzurePublicCloud":       "public",
	}

	for input, expected := range testData {
		actual := normalizeEnvironmentName(input)
		if actual != expected {
			t.Fatalf("Expected %q for input %q: got %q!", expected, input, actual)
		}
	}
}
