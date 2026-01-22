// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
)

var _ resourceids.Id = DnsCnameRecordPublicIpAddressAssociationId{}

type DnsCnameRecordPublicIpAddressAssociationId struct {
	CnameRecordId     recordsets.RecordTypeId
	PublicIpAddressId commonids.PublicIPAddressId
}

func NewDnsCnameRecordPublicIpAddressAssociationId(cnameRecordId recordsets.RecordTypeId, publicIpAddressId commonids.PublicIPAddressId) DnsCnameRecordPublicIpAddressAssociationId {
	return DnsCnameRecordPublicIpAddressAssociationId{
		CnameRecordId:     cnameRecordId,
		PublicIpAddressId: publicIpAddressId,
	}
}

func (id DnsCnameRecordPublicIpAddressAssociationId) ID() string {
	return fmt.Sprintf("%s|%s", id.CnameRecordId.ID(), id.PublicIpAddressId.ID())
}

func (id DnsCnameRecordPublicIpAddressAssociationId) String() string {
	components := []string{
		fmt.Sprintf("CnameRecordId %s", id.CnameRecordId.ID()),
		fmt.Sprintf("PublicIpAddressId %s", id.PublicIpAddressId.ID()),
	}
	return fmt.Sprintf("DNS CNAME Record Public IP Address Association: %s", strings.Join(components, " / "))
}

func DnsCnameRecordPublicIpAddressAssociationID(input string) (DnsCnameRecordPublicIpAddressAssociationId, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return DnsCnameRecordPublicIpAddressAssociationId{}, fmt.Errorf("expected ID to be in the format {CnameRecordId}|{PublicIpAddressId} but got %q", input)
	}

	cnameRecordId, err := recordsets.ParseRecordTypeID(parts[0])
	if err != nil {
		return DnsCnameRecordPublicIpAddressAssociationId{}, fmt.Errorf("parsing CNAME Record ID: %+v", err)
	}

	if cnameRecordId.RecordType != recordsets.RecordTypeCNAME {
		return DnsCnameRecordPublicIpAddressAssociationId{}, fmt.Errorf("expected record type to be CNAME but got %s", cnameRecordId.RecordType)
	}

	publicIpAddressId, err := commonids.ParsePublicIPAddressID(parts[1])
	if err != nil {
		return DnsCnameRecordPublicIpAddressAssociationId{}, fmt.Errorf("parsing Public IP Address ID: %+v", err)
	}

	if cnameRecordId == nil || publicIpAddressId == nil {
		return DnsCnameRecordPublicIpAddressAssociationId{}, fmt.Errorf("parse error, both CnameRecordId and PublicIpAddressId should not be nil")
	}

	return DnsCnameRecordPublicIpAddressAssociationId{
		CnameRecordId:     *cnameRecordId,
		PublicIpAddressId: *publicIpAddressId,
	}, nil
}

func DnsCnameRecordPublicIpAddressAssociationIDValidation(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := DnsCnameRecordPublicIpAddressAssociationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
