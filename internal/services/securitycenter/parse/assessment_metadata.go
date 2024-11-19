// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AssessmentMetadataId struct {
	SubscriptionId         string
	AssessmentMetadataName string
}

func NewAssessmentMetadataID(subscriptionId, assessmentMetadataName string) AssessmentMetadataId {
	return AssessmentMetadataId{
		SubscriptionId:         subscriptionId,
		AssessmentMetadataName: assessmentMetadataName,
	}
}

func (id AssessmentMetadataId) String() string {
	segments := []string{
		fmt.Sprintf("Assessment Metadata Name %q", id.AssessmentMetadataName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Assessment Metadata", segmentsStr)
}

func (id AssessmentMetadataId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/assessmentMetadata/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.AssessmentMetadataName)
}

// AssessmentMetadataID parses a AssessmentMetadata ID into an AssessmentMetadataId struct
func AssessmentMetadataID(input string) (*AssessmentMetadataId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AssessmentMetadata ID: %+v", input, err)
	}

	resourceId := AssessmentMetadataId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.AssessmentMetadataName, err = id.PopSegment("assessmentMetadata"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
