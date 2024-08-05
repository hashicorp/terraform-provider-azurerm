// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/privateendpoints"
)

var _ resourceids.Id = PrivateEndpointApplicationSecurityGroupAssociationId{}

type PrivateEndpointApplicationSecurityGroupAssociationId struct {
	PrivateEndpointId          privateendpoints.PrivateEndpointId
	ApplicationSecurityGroupId applicationsecuritygroups.ApplicationSecurityGroupId
}

func (p PrivateEndpointApplicationSecurityGroupAssociationId) ID() string {
	return fmt.Sprintf("%s|%s", p.PrivateEndpointId.ID(), p.ApplicationSecurityGroupId.ID())
}

func (p PrivateEndpointApplicationSecurityGroupAssociationId) String() string {
	components := []string{
		fmt.Sprintf("PrivateEndpointId %s", p.PrivateEndpointId.ID()),
		fmt.Sprintf("ApplicationSecurityGroupId %s", p.ApplicationSecurityGroupId.ID()),
	}
	return fmt.Sprintf("Private Endpoint Application Security Group Association: %s", strings.Join(components, " / "))
}

func NewPrivateEndpointApplicationSecurityGroupAssociationId(endpointId privateendpoints.PrivateEndpointId, securityGroupId applicationsecuritygroups.ApplicationSecurityGroupId) PrivateEndpointApplicationSecurityGroupAssociationId {
	return PrivateEndpointApplicationSecurityGroupAssociationId{
		PrivateEndpointId:          endpointId,
		ApplicationSecurityGroupId: securityGroupId,
	}
}

func PrivateEndpointApplicationSecurityGroupAssociationID(input string) (PrivateEndpointApplicationSecurityGroupAssociationId, error) {
	splitId := strings.Split(input, "|")
	if len(splitId) != 2 {
		return PrivateEndpointApplicationSecurityGroupAssociationId{}, fmt.Errorf("expected ID to be in the format {PrivateEndpointId}|{ApplicationSecurityGroupId} but got %q", input)
	}

	endpointId, err := privateendpoints.ParsePrivateEndpointID(splitId[0])
	if err != nil {
		return PrivateEndpointApplicationSecurityGroupAssociationId{}, err
	}

	securityGroupId, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(splitId[1])
	if err != nil {
		return PrivateEndpointApplicationSecurityGroupAssociationId{}, err
	}

	if endpointId == nil || securityGroupId == nil {
		return PrivateEndpointApplicationSecurityGroupAssociationId{}, fmt.Errorf("parse error, both PrivateEndpointId and ApplicationSecurityGroupId should not be nil")
	}

	return PrivateEndpointApplicationSecurityGroupAssociationId{
		PrivateEndpointId:          *endpointId,
		ApplicationSecurityGroupId: *securityGroupId,
	}, nil
}

func PrivateEndpointApplicationSecurityGroupAssociationIDValidation(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := PrivateEndpointApplicationSecurityGroupAssociationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
