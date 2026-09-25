// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sreagent

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2026-10-01/agents"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	sreclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sreagent/client"
)

const sreAgentExpectedStableAPIVersion = "2026-10-01"

// Server-provided continuation and polling URLs deliberately use a different version.
const sreAgentServerLinkAPIVersion = "2026-01-01"

func sreAgentVersionClient(t *testing.T, transport sreAgentTransport) (*agents.AgentsClient, agents.AgentId, context.Context) {
	t.Helper()
	c, err := sreclient.NewClient(&common.ClientOptions{
		Environment:                 *environments.AzurePublic(),
		Authorizers:                 &common.Authorizers{},
		Transport:                   transport,
		DisableCorrelationRequestID: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	c.Agents.Client.DisableRetries = true
	c.Agents.Client.AuthorizeRequest = nil
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)
	return c.Agents, agents.NewAgentID("00000000-0000-0000-0000-000000000000", "example", "agent"), ctx
}

func assertSreAgentRequestVersion(t *testing.T, request *http.Request, expected string) {
	t.Helper()
	if actual := request.URL.Query().Get("api-version"); actual != expected {
		t.Errorf("%s %s: api-version = %q, want %q", request.Method, request.URL.Path, actual, expected)
	}
}

func invokeSreAgentVersionOperation(ctx context.Context, c *agents.AgentsClient, id agents.AgentId, operation string) error {
	switch operation {
	case "read":
		_, err := c.Get(ctx, id)
		return err
	case "create":
		return c.CreateOrUpdateThenPoll(ctx, id, agents.Agent{
			Location: "eastus",
			Identity: &identity.LegacySystemAndUserAssignedMap{Type: identity.TypeSystemAssigned},
		})
	case "update":
		return updateSreAgent(ctx, c, id, agents.AgentPatch{Tags: pointer.To(map[string]string{"stage": "updated"})}, nil)
	case "identity_update":
		return updateSreAgent(ctx, c, id, agents.AgentPatch{
			Identity: &identity.LegacySystemAndUserAssignedMap{Type: identity.TypeSystemAssigned},
		}, nil)
	case "detach":
		return updateSreAgent(ctx, c, id, agents.AgentPatch{
			Properties: &agents.AgentPatchProperties{
				SandboxConfiguration: &agents.SandboxConfiguration{
					Egress: &agents.SandboxEgressConfiguration{Mode: pointer.To(agents.SandboxEgressModeLimited)},
				},
			},
		}, nil)
	case "delete":
		result, err := c.Delete(ctx, id)
		if err != nil {
			return err
		}
		return result.Poller.PollUntilDone(ctx)
	default:
		return fmt.Errorf("unknown test operation %q", operation)
	}
}

func TestSreAgentStableAPIVersion(t *testing.T) {
	for operation, method := range map[string]string{
		"read": http.MethodGet, "create": http.MethodPut, "update": http.MethodPatch,
		"identity_update": http.MethodPatch, "detach": http.MethodPatch, "delete": http.MethodDelete,
	} {
		t.Run(operation, func(t *testing.T) {
			operationCalls := 0
			deleted := false
			c, id, ctx := sreAgentVersionClient(t, func(request *http.Request) (*http.Response, error) {
				assertSreAgentRequestVersion(t, request, sreAgentExpectedStableAPIVersion)
				if request.Method == method {
					operationCalls++
				} else if request.Method != http.MethodGet {
					return nil, fmt.Errorf("unexpected method %s", request.Method)
				}
				if request.Method == http.MethodDelete {
					deleted = true
					return sreAgentHTTPResponse(request, http.StatusNoContent, ""), nil
				}
				if deleted {
					return sreAgentHTTPResponse(request, http.StatusNotFound, `{"error":{"code":"ResourceNotFound","message":"absent"}}`), nil
				}
				return sreAgentHTTPResponse(request, http.StatusOK, `{"location":"eastus","properties":{"provisioningState":"Succeeded"}}`), nil
			})
			if err := invokeSreAgentVersionOperation(ctx, c, id, operation); err != nil {
				t.Fatal(err)
			}
			if operationCalls != 1 {
				t.Fatalf("operation requests = %d, want 1", operationCalls)
			}
		})
	}
}

func TestSreAgentStableAPIVersionListContinuation(t *testing.T) {
	for _, scope := range []string{"subscription", "resource_group"} {
		t.Run(scope, func(t *testing.T) {
			path := "/subscriptions/00000000-0000-0000-0000-000000000000"
			if scope == "resource_group" {
				path += "/resourceGroups/example"
			}
			path += "/providers/Microsoft.App/agents"
			calls := 0
			c, id, ctx := sreAgentVersionClient(t, func(request *http.Request) (*http.Response, error) {
				calls++
				if request.Method != http.MethodGet || request.URL.Path != path {
					return nil, fmt.Errorf("unexpected List request %s %s", request.Method, request.URL)
				}
				if calls == 1 {
					assertSreAgentRequestVersion(t, request, sreAgentExpectedStableAPIVersion)
					body := fmt.Sprintf(`{"value":[],"nextLink":"https://management.azure.com%s?api-version=%s&$skiptoken=opaque%%2Ftoken"}`, path, sreAgentServerLinkAPIVersion)
					return sreAgentHTTPResponse(request, http.StatusOK, body), nil
				}
				if calls != 2 || request.URL.Query().Get("$skiptoken") != "opaque/token" {
					return nil, fmt.Errorf("unexpected continuation request %s", request.URL)
				}
				assertSreAgentRequestVersion(t, request, sreAgentServerLinkAPIVersion)
				return sreAgentHTTPResponse(request, http.StatusOK, `{"value":[]}`), nil
			})
			var err error
			if scope == "resource_group" {
				_, err = c.ListByResourceGroupComplete(ctx, commonids.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName))
			} else {
				_, err = c.ListBySubscriptionComplete(ctx, commonids.NewSubscriptionID(id.SubscriptionId))
			}
			if err != nil {
				t.Fatal(err)
			}
			if calls != 2 {
				t.Fatalf("List requests = %d, want 2", calls)
			}
		})
	}
}

func TestSreAgentStableAPIVersionPolling(t *testing.T) {
	for _, operation := range []string{"create", "update", "identity_update", "detach", "delete"} {
		t.Run(operation, func(t *testing.T) {
			writes, polls := 0, 0
			c, id, ctx := sreAgentVersionClient(t, func(request *http.Request) (*http.Response, error) {
				if request.URL.Path == "/operations/stable-version-test" {
					polls++
					if request.Method != http.MethodGet || request.URL.Query().Get("opaque") != "poll/token" {
						return nil, fmt.Errorf("unexpected polling request %s %s", request.Method, request.URL)
					}
					assertSreAgentRequestVersion(t, request, sreAgentServerLinkAPIVersion)
					return sreAgentHTTPResponse(request, http.StatusOK, `{"status":"Succeeded"}`), nil
				}
				assertSreAgentRequestVersion(t, request, sreAgentExpectedStableAPIVersion)
				if request.Method == http.MethodGet {
					if operation == "delete" {
						return sreAgentHTTPResponse(request, http.StatusNotFound, `{"error":{"code":"ResourceNotFound","message":"absent"}}`), nil
					}
					return sreAgentHTTPResponse(request, http.StatusOK, `{"properties":{"provisioningState":"Succeeded"}}`), nil
				}
				writes++
				status := http.StatusAccepted
				if request.Method == http.MethodPut {
					status = http.StatusCreated
				}
				response := sreAgentHTTPResponse(request, status, `{"properties":{"provisioningState":"Accepted"}}`)
				response.Header.Set("Azure-AsyncOperation", "https://management.azure.com/operations/stable-version-test?api-version="+sreAgentServerLinkAPIVersion+"&opaque=poll%2Ftoken")
				return response, nil
			})
			if err := invokeSreAgentVersionOperation(ctx, c, id, operation); err != nil {
				t.Fatal(err)
			}
			if writes != 1 || polls != 1 {
				t.Fatalf("writes=%d polls=%d, want one of each", writes, polls)
			}
		})
	}
}
