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

type JobId struct {
	SubscriptionId   string
	ResourceGroup    string
	BatchAccountName string
	PoolName         string
	Name             string
}

func NewJobID(subscriptionId, resourceGroup, batchAccountName, poolName, name string) JobId {
	return JobId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		BatchAccountName: batchAccountName,
		PoolName:         poolName,
		Name:             name,
	}
}

func (id JobId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Pool Name %q", id.PoolName),
		fmt.Sprintf("Batch Account Name %q", id.BatchAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Job", segmentsStr)
}

func (id JobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Batch/batchAccounts/%s/pools/%s/jobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BatchAccountName, id.PoolName, id.Name)
}

// JobID parses a Job ID into an JobId struct
func JobID(input string) (*JobId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Job ID: %+v", input, err)
	}

	resourceId := JobId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.BatchAccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}
	if resourceId.PoolName, err = id.PopSegment("pools"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("jobs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
