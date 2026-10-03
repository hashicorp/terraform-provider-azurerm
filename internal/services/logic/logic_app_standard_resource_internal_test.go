// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetBasicLogicAppSettings(t *testing.T) {
	testCases := []struct {
		name                 string
		secretID             string
		shareName            string
		useExtensionBundle   bool
		bundleVersion        string
		expectedConnection   string
		expectedContentShare string
		expectedError        string
	}{
		{
			name:                 "access key",
			expectedConnection:   "DefaultEndpointsProtocol=https;AccountName=teststorage;AccountKey=testkey;EndpointSuffix=core.windows.net",
			expectedContentShare: "testlogicapp-content",
		},
		{
			name:                 "key vault",
			secretID:             "https://test.vault.azure.net/secrets/storage",
			expectedConnection:   "@Microsoft.KeyVault(SecretUri=https://test.vault.azure.net/secrets/storage)",
			expectedContentShare: "testlogicapp-content",
		},
		{
			name:                 "custom share and extension bundle",
			shareName:            "custom-content",
			useExtensionBundle:   true,
			bundleVersion:        "[1.*, 2.0.0)",
			expectedConnection:   "DefaultEndpointsProtocol=https;AccountName=teststorage;AccountKey=testkey;EndpointSuffix=core.windows.net",
			expectedContentShare: "custom-content",
		},
		{
			name:               "missing bundle version",
			useExtensionBundle: true,
			expectedError:      "when `use_extension_bundle` is true, `bundle_version` must be specified",
		},
	}

	for _, testCase := range testCases {
		for _, isASE := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/ASE=%t", testCase.name, isASE), func(t *testing.T) {
				model := LogicAppResourceModel{
					Name:                    "TestLogicApp",
					StorageAccountName:      "teststorage",
					StorageAccountAccessKey: "testkey",
					StorageKeyVaultSecretID: testCase.secretID,
					StorageAccountShareName: testCase.shareName,
					Version:                 "~4",
					UseExtensionBundle:      testCase.useExtensionBundle,
					BundleVersion:           testCase.bundleVersion,
				}

				settings, err := getBasicLogicAppSettings(model, "core.windows.net", isASE)
				if testCase.expectedError != "" {
					if err == nil || err.Error() != testCase.expectedError {
						t.Fatalf("expected error %q, got %v", testCase.expectedError, err)
					}
					if settings != nil {
						t.Fatalf("expected no settings on error, got %#v", settings)
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}

				expected := map[string]string{
					"AzureWebJobsStorage":         testCase.expectedConnection,
					"FUNCTIONS_EXTENSION_VERSION": "~4",
					"APP_KIND":                    "workflowApp",
				}
				if !isASE {
					expected["WEBSITE_CONTENTSHARE"] = testCase.expectedContentShare
					expected["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"] = testCase.expectedConnection
				}
				if testCase.useExtensionBundle {
					expected["AzureFunctionsJobHost__extensionBundle__id"] = "Microsoft.Azure.Functions.ExtensionBundle.Workflows"
					expected["AzureFunctionsJobHost__extensionBundle__version"] = testCase.bundleVersion
				}

				actual := make(map[string]string, len(settings))
				for _, setting := range settings {
					if setting.Name == nil || setting.Value == nil {
						t.Fatalf("expected a non-nil setting name and value, got %#v", setting)
					}
					if _, exists := actual[*setting.Name]; exists {
						t.Fatalf("duplicate app setting %q", *setting.Name)
					}
					actual[*setting.Name] = *setting.Value
				}
				if !reflect.DeepEqual(actual, expected) {
					t.Fatalf("expected settings %#v, got %#v", expected, actual)
				}
			})
		}
	}
}
