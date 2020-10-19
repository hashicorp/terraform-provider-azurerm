package parse

import (
	"fmt"
	"strings"
)

type ManagementGroupSubscriptionAssociationId struct {
	ManagementGroupID   string
	ManagementGroupName string
	SubscriptionScopeID string
	SubscriptionID      string
}

func ManagementGroupSubscriptionAssociationID(input string) (*ManagementGroupSubscriptionAssociationId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected an ID in the format `{managementGroupId}|{subscriptionId} but got %q", input)
	}

	managementGroupId := segments[0]
	subscriptionScopeId := segments[1]

	parsedManagementGroupID, err := ManagementGroupID(managementGroupId)
	if err != nil {
		return nil, fmt.Errorf("parsing Management Group ID %q: %+v", managementGroupId, err)
	}
	managementGroupName := parsedManagementGroupID.Name

	// The subscriptionId is the full scope form, i.e. '/subscriptions/{subscriptionGuid}'
	if subscriptionScopeSegments := strings.Split(subscriptionScopeId, "/"); len(subscriptionScopeSegments) != 3 {
		return nil, fmt.Errorf("unable to parse subscription scope ID %q - has %d segments", subscriptionScopeId, len(subscriptionScopeSegments))
	}

	subscriptionId := strings.TrimPrefix(subscriptionScopeId, "/subscriptions/")
	if subscriptionId == subscriptionScopeId {
		return nil, fmt.Errorf("expected subscription scope ID %q  in `/subscriptions/00000000-0000-0000-0000-000000000000` format", subscriptionScopeId)
	}

	id := ManagementGroupSubscriptionAssociationId{
		ManagementGroupID:   managementGroupId,
		ManagementGroupName: managementGroupName, // strings.Split(managementGroupId, "/")[-1],
		SubscriptionScopeID: subscriptionScopeId,
		SubscriptionID:      subscriptionId,
	}

	return &id, nil
}
