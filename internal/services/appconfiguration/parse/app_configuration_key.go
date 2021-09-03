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
	cfgStoreId := strings.TrimSuffix(input, fmt.Sprintf("/AppConfigurationKey/%s/Label/%s", keyName, label))
	appcfgID := AppConfigurationKeyId{
		ConfigurationStoreId: cfgStoreId,
		Key:                  keyName,
		Label:                label,
	}

	return &appcfgID, nil
}
