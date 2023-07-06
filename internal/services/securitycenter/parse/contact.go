// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ContactId struct {
	SubscriptionId      string
	SecurityContactName string
}

func NewContactID(subscriptionId, securityContactName string) ContactId {
	return ContactId{
		SubscriptionId:      subscriptionId,
		SecurityContactName: securityContactName,
	}
}

func (id ContactId) String() string {
	segments := []string{
		fmt.Sprintf("Security Contact Name %q", id.SecurityContactName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Contact", segmentsStr)
}

func (id ContactId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/securityContacts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.SecurityContactName)
}

// ContactID parses a Contact ID into an ContactId struct
func ContactID(input string) (*ContactId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Contact ID: %+v", input, err)
	}

	resourceId := ContactId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.SecurityContactName, err = id.PopSegment("securityContacts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
