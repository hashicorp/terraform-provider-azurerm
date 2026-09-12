// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ServicePrincipalAttributeSetId struct {
	ServicePrincipalObjectId string
	AttributeSetName         string
}

func NewServicePrincipalAttributeSetId(servicePrincipalObjectId string, attributeSetName string) ServicePrincipalAttributeSetId {
	return ServicePrincipalAttributeSetId{
		ServicePrincipalObjectId: servicePrincipalObjectId,
		AttributeSetName:         attributeSetName,
	}
}

func (id ServicePrincipalAttributeSetId) ID() string {
	return fmt.Sprintf("%s|%s", id.ServicePrincipalObjectId, id.AttributeSetName)
}

func (id ServicePrincipalAttributeSetId) String() string {
	return fmt.Sprintf("Service Principal Object ID %q / Attribute Set %q", id.ServicePrincipalObjectId, id.AttributeSetName)
}

func ServicePrincipalAttributeSetID(input string) (*ServicePrincipalAttributeSetId, error) {
	parts := strings.SplitN(input, "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("parsing Service Principal Attribute Set ID: expected '<service_principal_object_id>|<attribute_set_name>' but got %q", input)
	}

	if _, errs := validation.IsUUID(parts[0], "service_principal_object_id"); len(errs) > 0 {
		return nil, fmt.Errorf("parsing Service Principal Attribute Set ID: invalid service principal object ID %q", parts[0])
	}

	if parts[1] == "" {
		return nil, fmt.Errorf("parsing Service Principal Attribute Set ID: attribute set name cannot be empty")
	}

	return &ServicePrincipalAttributeSetId{
		ServicePrincipalObjectId: parts[0],
		AttributeSetName:         parts[1],
	}, nil
}

func ValidateServicePrincipalAttributeSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ServicePrincipalAttributeSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
