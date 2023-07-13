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

var _ resourceids.Id = AppConfigurationFeatureId{}

type AppConfigurationFeatureId struct {
	ConfigurationStoreId string
	Name                 string
	Label                string
}

func (id AppConfigurationFeatureId) ID() string {
	return fmt.Sprintf("%s/AppConfigurationFeature/%s/Label/%s", id.ConfigurationStoreId, id.Name, id.Label)
}

func (id AppConfigurationFeatureId) String() string {
	components := []string{
		fmt.Sprintf("Configuration Store Id %q", id.ConfigurationStoreId),
		fmt.Sprintf("Label %q", id.Label),
		fmt.Sprintf("Name %q", id.Name),
	}
	return fmt.Sprintf("Feature: %s", strings.Join(components, " / "))
}

func FeatureId(input string) (*AppConfigurationFeatureId, error) {
	resourceID, err := parseAzureResourceID(handleSlashInIdForFeature(input))
	if err != nil {
		return nil, fmt.Errorf("while parsing resource ID: %+v", err)
	}

	keyName := resourceID.Path["AppConfigurationFeature"]
	label := resourceID.Path["Label"]

	appcfgID := AppConfigurationFeatureId{
		Name:  keyName,
		Label: label,
	}

	// Golang's URL parser will translate %00 to \000 (NUL). This will only happen if we're dealing with an empty
	// label, so we set the label to the expected value (empty string) and trim the input string, so we can properly
	// extract the configuration store ID out of it.
	if label == "\000" {
		appcfgID.Label = ""
		input = strings.TrimSuffix(input, "%00")
	}
	appcfgID.ConfigurationStoreId = strings.TrimSuffix(input, fmt.Sprintf("/AppConfigurationFeature/%s/Label/%s", appcfgID.Name, appcfgID.Label))

	return &appcfgID, nil
}

// a workaround to support "/" in id
func handleSlashInIdForFeature(input string) string {
	oldNames := regexp.MustCompile(`AppConfigurationFeature\/(.+)\/Label\/`).FindStringSubmatch(input)
	if len(oldNames) == 2 {
		input = strings.Replace(input, oldNames[1], url.QueryEscape(oldNames[1]), 1)
	}

	oldNames = regexp.MustCompile(`AppConfigurationFeature\/.+\/Label\/(.+)`).FindStringSubmatch(input)

	// Label will have a "%00" placeholder if we're dealing with an empty label,
	if len(oldNames) == 2 && oldNames[1] != "%00" {
		input = strings.Replace(input, oldNames[1], url.QueryEscape(oldNames[1]), 1)
	}

	return input
}
