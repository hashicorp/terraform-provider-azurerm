// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = AppConfigurationKeyId{}

type AppConfigurationKeyId struct {
	ConfigurationStoreId string
	Key                  string
	Label                string
}

func (id AppConfigurationKeyId) ID() string {
	return fmt.Sprintf("%s/AppConfigurationKey/%s/Label/%s", id.ConfigurationStoreId, id.Key, id.Label)
}

func (id AppConfigurationKeyId) String() string {
	components := []string{
		fmt.Sprintf("Configuration Store Id %q", id.ConfigurationStoreId),
		fmt.Sprintf("Label %q", id.Label),
		fmt.Sprintf("Key %q", id.Key),
	}
	return fmt.Sprintf("Key: %s", strings.Join(components, " / "))
}

func KeyId(input string) (*AppConfigurationKeyId, error) {
	resourceID, err := parseAzureResourceID(handleSlashInIdForKey(input))
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

// a workaround to support "/" in id
func handleSlashInIdForKey(input string) string {
	oldNames := regexp.MustCompile(`AppConfigurationKey\/(.+)\/Label\/`).FindStringSubmatch(input)
	if len(oldNames) == 2 {
		input = strings.Replace(input, oldNames[1], url.QueryEscape(oldNames[1]), 1)
	}

	oldNames = regexp.MustCompile(`AppConfigurationKey\/.+\/Label\/(.+)`).FindStringSubmatch(input)

	// Label will have a "%00" placeholder if we're dealing with an empty label,
	if len(oldNames) == 2 && oldNames[1] != "%00" {
		input = strings.Replace(input, oldNames[1], url.QueryEscape(oldNames[1]), 1)
	}

	return input

}
