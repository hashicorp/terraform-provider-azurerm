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

type JobCredentialId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerName     string
	JobAgentName   string
	CredentialName string
}

func NewJobCredentialID(subscriptionId, resourceGroup, serverName, jobAgentName, credentialName string) JobCredentialId {
	return JobCredentialId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerName:     serverName,
		JobAgentName:   jobAgentName,
		CredentialName: credentialName,
	}
}

func (id JobCredentialId) String() string {
	segments := []string{
		fmt.Sprintf("Credential Name %q", id.CredentialName),
		fmt.Sprintf("Job Agent Name %q", id.JobAgentName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Job Credential", segmentsStr)
}

func (id JobCredentialId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/jobAgents/%s/credentials/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.JobAgentName, id.CredentialName)
}

// JobCredentialID parses a JobCredential ID into an JobCredentialId struct
func JobCredentialID(input string) (*JobCredentialId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an JobCredential ID: %+v", input, err)
	}

	resourceId := JobCredentialId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.JobAgentName, err = id.PopSegment("jobAgents"); err != nil {
		return nil, err
	}
	if resourceId.CredentialName, err = id.PopSegment("credentials"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
