// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NestedItemObjectType string

const (
	NestedItemTypeOrganization NestedItemObjectType = "organizations"
)

// PossibleNestedItemObjectTypeValues returns a string slice of possible "NestedItemObjectType" values.
func PossibleNestedItemObjectTypeValues() []string {
	return []string{string(NestedItemTypeOrganization)}
}

var _ resourceids.Id = NestedItemId{}

type NestedItemId struct {
	IotcentralBaseUrl string
	SubDomain         string
	NestedItemType    NestedItemObjectType
	Id                string
}

func NewNestedItemID(iotcentralBaseUrl string, nestedItemType NestedItemObjectType, id string) (*NestedItemId, error) {
	iotCentralUrl, err := url.Parse(iotcentralBaseUrl)
	if err != nil || iotcentralBaseUrl == "" {
		return nil, fmt.Errorf("parsing %q: %+v", iotcentralBaseUrl, err)
	}

	// in case the user has provided a port, we need to remove it
	if hostParts := strings.Split(iotCentralUrl.Host, ":"); len(hostParts) > 1 {
		iotCentralUrl.Host = hostParts[0]
	}

	subDomain, err := parseSubDomain(iotcentralBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", iotcentralBaseUrl, err)
	}

	return &NestedItemId{
		IotcentralBaseUrl: iotCentralUrl.String(),
		SubDomain:         subDomain,
		NestedItemType:    nestedItemType,
		Id:                id,
	}, nil
}

func (id NestedItemId) ID() string {
	// example: https://subdomain.baseDomain/api/organizations/fdf067c93bbb4b22bff4d8b7a9a56217
	segments := []string{
		strings.TrimSuffix(id.IotcentralBaseUrl, "/"),
		"api",
		string(id.NestedItemType),
		id.Id,
	}

	return strings.TrimSuffix(strings.Join(segments, "/"), "/")
}

func (id NestedItemId) String() string {
	components := []string{
		fmt.Sprintf("Base Url %q", id.IotcentralBaseUrl),
		fmt.Sprintf("Sub Domain %q", id.SubDomain),
		fmt.Sprintf("Nested Item Type %q", string(id.NestedItemType)),
		fmt.Sprintf("Id %q", id.Id),
	}
	return fmt.Sprintf("Iot Central Nested Item %s", strings.Join(components, " / "))
}

func ParseNestedItemID(input string) (*NestedItemId, error) {
	item, err := parseNestedItemId(input)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func parseNestedItemId(id string) (*NestedItemId, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("cannot parse azure iot central child ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	if len(components) != 3 {
		return nil, fmt.Errorf("iot central nested item should 3 segments, found %d segment(s) in %q", len(components), id)
	}

	nestedItemObjectTypes := PossibleNestedItemObjectTypeValues()

	if !utils.SliceContainsValue(nestedItemObjectTypes, components[1]) {
		return nil, fmt.Errorf("key vault 'NestedItemType' should be one of: %s, got %q", strings.Join(nestedItemObjectTypes, ", "), components[1])
	}

	nestedItemObjectType := NestedItemObjectType(components[1])

	subdomain, err := parseSubDomain(idURL.Scheme + "://" + idURL.Host)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", idURL.Scheme+"://"+idURL.Host, err)
	}

	childId := NestedItemId{
		IotcentralBaseUrl: fmt.Sprintf("%s://%s/", idURL.Scheme, idURL.Host),
		SubDomain:         subdomain,
		NestedItemType:    nestedItemObjectType,
		Id:                components[2],
	}

	return &childId, nil
}

func parseSubDomain(iotcentralBaseUrl string) (string, error) {
	iotcentralUrl, err := url.Parse(iotcentralBaseUrl)
	if err != nil || iotcentralBaseUrl == "" {
		return "", fmt.Errorf("parsing %q: %+v", iotcentralBaseUrl, err)
	}

	subDomain := strings.Split(iotcentralUrl.Host, ".")[0]

	return subDomain, nil
}
