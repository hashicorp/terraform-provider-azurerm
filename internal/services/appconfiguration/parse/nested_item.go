// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = NestedItemId{}

type NestedItemId struct {
	ConfigurationStoreEndpoint string
	Key                        string
	Label                      string
}

func NewNestedItemID(configurationStoreEndpoint, key, label string) (*NestedItemId, error) {
	// configurationStoreEndpoint example: https://testappconf1.azconfig.io
	configurationURL, err := url.ParseRequestURI(configurationStoreEndpoint)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", configurationURL, err)
	}

	return &NestedItemId{
		ConfigurationStoreEndpoint: configurationURL.String(),
		Key:                        key,
		Label:                      label,
	}, nil
}

func (id NestedItemId) ID() string {
	// example with label: https://testappconf1.azconfig.io/kv/testKey?label=testLabel
	// example without label: https://testappconf1.azconfig.io/kv/testKey?laebl=
	baseURL, _ := url.ParseRequestURI(id.ConfigurationStoreEndpoint)
	u := &url.URL{
		Scheme:   baseURL.Scheme,
		Host:     baseURL.Host,
		Path:     fmt.Sprintf("kv/%s", id.Key),
		RawPath:  fmt.Sprintf("kv/%s", url.PathEscape(id.Key)),
		RawQuery: fmt.Sprintf("label=%s", url.QueryEscape(id.Label)),
	}

	return u.String()
}

func (id NestedItemId) String() string {
	components := []string{
		fmt.Sprintf("Configuration Store Endpoint %q", id.ConfigurationStoreEndpoint),
		fmt.Sprintf("Key %q", id.Key),
		fmt.Sprintf("Label %q", id.Label),
	}
	return fmt.Sprintf("AppConfiguration Nested Item %s", strings.Join(components, " / "))
}

// ParseNestedItemID parses an App Configuration Nested Item ID (such as a Key or Feature)
func ParseNestedItemID(input string) (*NestedItemId, error) {
	// example with label: https://testappconf1.azconfig.io/kv/testKey?label=testLabel
	// example without label: https://testappconf1.azconfig.io/kv/testKey?label=
	idURL, err := url.ParseRequestURI(input)
	if err != nil {
		return nil, fmt.Errorf("cannot parse Azure App Configuration Nested Item ID %q: %s", input, err)
	}

	rawPath := idURL.EscapedPath()
	rawPath = strings.TrimPrefix(rawPath, "/")
	rawPath = strings.TrimSuffix(rawPath, "/")

	components := strings.Split(rawPath, "/")
	if len(components) != 2 {
		return nil, fmt.Errorf("AppConfiguration Nested Item should contain 2 segments, got %d from %q", len(components), rawPath)
	}

	key, err := url.PathUnescape(components[1])
	if err != nil {
		return nil, fmt.Errorf("cannot unescape Azure App Configuration Nested Item key %q: %s", components[1], err)
	}

	label := ""
	queryMap := idURL.Query()
	rawLabel, ok := queryMap["label"]

	if len(queryMap) != 1 || !ok || len(rawLabel) != 1 {
		return nil, fmt.Errorf("exactly one 'label' must be defined in Azure App Configuration Nested Item URL query, but got %q", idURL.RawQuery)
	}

	// Golang's URL parser will translate %00 to \000 (NUL). This will only happen if we're dealing with an empty
	// label, so we set the label to the expected value (empty string)
	if rawLabel[0] != "\000" {
		label = rawLabel[0]
	}

	if label == "" && !strings.HasSuffix(input, "label=") {
		return nil, fmt.Errorf("only 'label=' is accepted for empty label in Azure App Configuration Nested Item URL query")
	}

	childId := NestedItemId{
		ConfigurationStoreEndpoint: fmt.Sprintf("%s://%s", idURL.Scheme, idURL.Host),
		Key:                        key,
		Label:                      label,
	}

	return &childId, nil
}
