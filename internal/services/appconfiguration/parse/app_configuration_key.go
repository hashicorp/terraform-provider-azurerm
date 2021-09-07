package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AppConfigurationKeyId struct {
	ConfigurationStoreId string
	Key                  string
	Label                string
}

func (k AppConfigurationKeyId) ID() string {
	return fmt.Sprintf("%s/AppConfigurationKey/%s/Label/%s", k.ConfigurationStoreId, k.Key, k.Label)
}

func AppConfigurationKeyID(input string) (*AppConfigurationKeyId, error) {

	resourceID, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	keyName := resourceID.Path["AppConfigurationKey"]
	label := resourceID.Path["Label"]

	appcfgID := AppConfigurationKeyId{
		Key:   keyName,
		Label: label,
	}

	// Golang's URL parser will translate %00 to \000 (NUL). This will only happen if we're dealing with an empty
	// label, so we set the label to the expected value (empty string) and trim the input string, so we can properly
	// extract the configuration store ID out of it.
	if label == "\000" {
		appcfgID.Label = ""
		input = strings.TrimSuffix(input, "%00")
	}
	appcfgID.ConfigurationStoreId = strings.TrimSuffix(input, fmt.Sprintf("/AppConfigurationKey/%s/Label/%s", appcfgID.Key, appcfgID.Label))

	return &appcfgID, nil
}
