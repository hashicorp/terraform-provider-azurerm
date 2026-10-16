// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-10-01/agents"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
)

func validateSreAgentIdentity(value *identity.LegacySystemAndUserAssignedMap) error {
	if value == nil {
		return fmt.Errorf("`identity` must assign at least one managed identity")
	}
	switch value.Type {
	case identity.TypeSystemAssigned:
		return nil
	case identity.TypeUserAssigned, identity.TypeSystemAssignedUserAssigned:
		if len(value.IdentityIds) > 0 {
			return nil
		}
		return fmt.Errorf("`identity_ids` must not be empty when `identity.type` includes UserAssigned")
	default:
		return fmt.Errorf("`identity` must assign at least one managed identity; removing the last identity is not supported")
	}
}

type sreAgentIdentityPatch struct {
	Type                   string                                           `json:"type"`
	UserAssignedIdentities map[string]*identity.UserAssignedIdentityDetails `json:"userAssignedIdentities"`
}

func expandSreAgentIdentityPatch(desired, previous *identity.LegacySystemAndUserAssignedMap) *sreAgentIdentityPatch {
	result := &sreAgentIdentityPatch{Type: strings.ReplaceAll(string(desired.Type), ", ", ",")}
	if desired.Type == identity.TypeSystemAssigned {
		return result
	}
	result.UserAssignedIdentities = make(map[string]*identity.UserAssignedIdentityDetails)
	retained := make(map[string]struct{}, len(desired.IdentityIds))
	for id := range desired.IdentityIds {
		result.UserAssignedIdentities[id] = &identity.UserAssignedIdentityDetails{}
		retained[strings.ToLower(id)] = struct{}{}
	}
	if previous != nil {
		for id := range previous.IdentityIds {
			if _, exists := retained[strings.ToLower(id)]; !exists {
				result.UserAssignedIdentities[id] = nil
			}
		}
	}
	return result
}

func updateSreAgent(ctx context.Context, c *agents.AgentsClient, id agents.AgentId, input agents.AgentPatch, previous *identity.LegacySystemAndUserAssignedMap) error {
	if input.Identity == nil && !sreAgentNeedsNetworkDetach(input.Properties) {
		return c.UpdateThenPoll(ctx, id, input)
	}
	var identityPatch *sreAgentIdentityPatch
	if input.Identity != nil {
		if err := validateSreAgentIdentity(input.Identity); err != nil {
			return err
		}
		identityPatch = expandSreAgentIdentityPatch(input.Identity, previous)
	}
	properties, err := marshalSreAgentPatchProperties(input.Properties)
	if err != nil {
		return fmt.Errorf("encoding selective SRE Agent patch properties: %+v", err)
	}
	// Compose identity tombstones and explicit network detach without changing other owned fields.
	payload := struct {
		agents.AgentPatch
		Identity   *sreAgentIdentityPatch `json:"identity,omitempty"`
		Properties json.RawMessage        `json:"properties,omitempty"`
	}{
		AgentPatch: input,
		Identity:   identityPatch,
		Properties: properties,
	}
	req, err := c.Client.NewRequest(ctx, client.RequestOptions{
		ContentType:         "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{http.StatusOK, http.StatusAccepted},
		HttpMethod:          http.MethodPatch,
		Path:                id.ID(),
	})
	if err != nil {
		return err
	}
	if err := req.Marshal(payload); err != nil {
		return err
	}
	resp, err := req.Execute(ctx)
	if err != nil {
		return err
	}
	poller, err := resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return err
	}
	return poller.PollUntilDone(ctx)
}
