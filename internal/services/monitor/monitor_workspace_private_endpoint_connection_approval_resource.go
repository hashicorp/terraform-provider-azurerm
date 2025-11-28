// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkspacePrivateEndpointConnectionResourceModel struct {
	Name                          string `tfschema:"name"`
	WorkspaceId                   string `tfschema:"workspace_id"`
	PrivateEndpointConnectionName string `tfschema:"private_endpoint_connection_name"`
	ApprovalMessage               string `tfschema:"approval_message"`
	PrivateEndpointId             string `tfschema:"private_endpoint_id"`
}

type WorkspacePrivateEndpointConnectionResource struct{}

var _ sdk.Resource = WorkspacePrivateEndpointConnectionResource{}

func (r WorkspacePrivateEndpointConnectionResource) ResourceType() string {
	return "azurerm_monitor_workspace_private_endpoint_connection_approval"
}

func (r WorkspacePrivateEndpointConnectionResource) ModelObject() interface{} {
	return &WorkspacePrivateEndpointConnectionResourceModel{}
}

func (r WorkspacePrivateEndpointConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return azuremonitorworkspaces.ValidateAccountID
}

func (r WorkspacePrivateEndpointConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azuremonitorworkspaces.ValidateAccountID,
		},

		"private_endpoint_connection_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"approval_message": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Approved via Terraform",
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r WorkspacePrivateEndpointConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_endpoint_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r WorkspacePrivateEndpointConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspacePrivateEndpointConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceClient := metadata.Client.Monitor.WorkspacesClient

			workspaceId, err := azuremonitorworkspaces.ParseAccountID(model.WorkspaceId)
			if err != nil {
				return err
			}

			// Wait for the private endpoint connection to appear (handles eventual consistency)
			// This is needed because Azure may take some time to show the connection after
			// a managed private endpoint is created
			targetConnection, err := waitForPrivateEndpointConnectionToAppear(ctx, workspaceClient, *workspaceId, model.PrivateEndpointConnectionName)
			if err != nil {
				return fmt.Errorf("waiting for private endpoint connection %q on %s: %+v", model.PrivateEndpointConnectionName, *workspaceId, err)
			}

			// Check if already approved (requires import)
			if targetConnection.Properties != nil &&
				targetConnection.Properties.PrivateLinkServiceConnectionState.Status != nil &&
				*targetConnection.Properties.PrivateLinkServiceConnectionState.Status == azuremonitorworkspaces.PrivateEndpointServiceConnectionStatusApproved {
				return metadata.ResourceRequiresImport(r.ResourceType(), workspaceId)
			}

			// Build the approval request body
			approvalBody := azuremonitorworkspaces.PrivateEndpointConnection{
				Properties: &azuremonitorworkspaces.PrivateEndpointConnectionProperties{
					PrivateLinkServiceConnectionState: azuremonitorworkspaces.PrivateLinkServiceConnectionState{
						Status:          pointer.To(azuremonitorworkspaces.PrivateEndpointServiceConnectionStatusApproved),
						Description:     pointer.To(model.ApprovalMessage),
						ActionsRequired: pointer.To("None"),
					},
				},
			}

			// Make the PUT request to approve the connection
			connectionPath := fmt.Sprintf("%s/privateEndpointConnections/%s", workspaceId.ID(), model.PrivateEndpointConnectionName)

			opts := client.RequestOptions{
				ContentType: "application/json; charset=utf-8",
				ExpectedStatusCodes: []int{
					http.StatusOK,
					http.StatusCreated,
					http.StatusAccepted,
				},
				HttpMethod: http.MethodPut,
				Path:       connectionPath,
			}

			req, err := workspaceClient.Client.NewRequest(ctx, opts)
			if err != nil {
				return fmt.Errorf("building request: %+v", err)
			}

			if err := req.Marshal(&approvalBody); err != nil {
				return fmt.Errorf("marshaling request body: %+v", err)
			}

			resp, err := workspaceClient.Client.Execute(ctx, req)
			if err != nil {
				return fmt.Errorf("approving private endpoint connection %q on %s: %+v", model.PrivateEndpointConnectionName, *workspaceId, err)
			}

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
				return fmt.Errorf("unexpected status code %d when approving private endpoint connection", resp.StatusCode)
			}

			// Use the workspace ID as the resource ID since the connection is a sub-resource
			metadata.SetID(workspaceId)

			return nil
		},
	}
}

func (r WorkspacePrivateEndpointConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			workspaceClient := metadata.Client.Monitor.WorkspacesClient

			workspaceId, err := azuremonitorworkspaces.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Get stored connection name from state
			var model WorkspacePrivateEndpointConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspace, err := workspaceClient.Get(ctx, *workspaceId)
			if err != nil {
				if response.WasNotFound(workspace.HttpResponse) {
					return metadata.MarkAsGone(workspaceId)
				}
				return fmt.Errorf("retrieving %s: %+v", *workspaceId, err)
			}

			if workspace.Model == nil || workspace.Model.Properties == nil || workspace.Model.Properties.PrivateEndpointConnections == nil {
				return metadata.MarkAsGone(workspaceId)
			}

			// Find the connection by name
			var targetConnection *azuremonitorworkspaces.PrivateEndpointConnection
			for _, conn := range *workspace.Model.Properties.PrivateEndpointConnections {
				if conn.Name != nil && *conn.Name == model.PrivateEndpointConnectionName {
					targetConnection = &conn
					break
				}
			}

			if targetConnection == nil {
				return metadata.MarkAsGone(workspaceId)
			}

			// Check if it's still approved
			if targetConnection.Properties != nil &&
				targetConnection.Properties.PrivateLinkServiceConnectionState.Status != nil &&
				*targetConnection.Properties.PrivateLinkServiceConnectionState.Status != azuremonitorworkspaces.PrivateEndpointServiceConnectionStatusApproved {
				return metadata.MarkAsGone(workspaceId)
			}

			state := WorkspacePrivateEndpointConnectionResourceModel{
				WorkspaceId:                   workspaceId.ID(),
				PrivateEndpointConnectionName: model.PrivateEndpointConnectionName,
				ApprovalMessage:               model.ApprovalMessage,
				Name:                          pointer.From(targetConnection.Name),
			}

			if targetConnection.Properties != nil && targetConnection.Properties.PrivateEndpoint != nil {
				state.PrivateEndpointId = pointer.From(targetConnection.Properties.PrivateEndpoint.Id)
			}

			if targetConnection.Properties != nil && targetConnection.Properties.PrivateLinkServiceConnectionState.Description != nil {
				state.ApprovalMessage = pointer.From(targetConnection.Properties.PrivateLinkServiceConnectionState.Description)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkspacePrivateEndpointConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspacePrivateEndpointConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceClient := metadata.Client.Monitor.WorkspacesClient

			workspaceId, err := azuremonitorworkspaces.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("approval_message") {
				// Re-approve with new message
				approvalBody := azuremonitorworkspaces.PrivateEndpointConnection{
					Properties: &azuremonitorworkspaces.PrivateEndpointConnectionProperties{
						PrivateLinkServiceConnectionState: azuremonitorworkspaces.PrivateLinkServiceConnectionState{
							Status:          pointer.To(azuremonitorworkspaces.PrivateEndpointServiceConnectionStatusApproved),
							Description:     pointer.To(model.ApprovalMessage),
							ActionsRequired: pointer.To("None"),
						},
					},
				}

				connectionPath := fmt.Sprintf("%s/privateEndpointConnections/%s", workspaceId.ID(), model.PrivateEndpointConnectionName)

				opts := client.RequestOptions{
					ContentType: "application/json; charset=utf-8",
					ExpectedStatusCodes: []int{
						http.StatusOK,
						http.StatusCreated,
						http.StatusAccepted,
					},
					HttpMethod: http.MethodPut,
					Path:       connectionPath,
				}

				req, err := workspaceClient.Client.NewRequest(ctx, opts)
				if err != nil {
					return fmt.Errorf("building request: %+v", err)
				}

				if err := req.Marshal(&approvalBody); err != nil {
					return fmt.Errorf("marshaling request body: %+v", err)
				}

				if _, err := workspaceClient.Client.Execute(ctx, req); err != nil {
					return fmt.Errorf("updating private endpoint connection %q on %s: %+v", model.PrivateEndpointConnectionName, *workspaceId, err)
				}
			}

			return nil
		},
	}
}

func (r WorkspacePrivateEndpointConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// On delete, we reject the connection (or just remove from state)
			// For now, we'll just remove from state without rejecting
			// as the connection will remain but be managed outside of Terraform

			var model WorkspacePrivateEndpointConnectionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceClient := metadata.Client.Monitor.WorkspacesClient

			workspaceId, err := azuremonitorworkspaces.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Reject the connection on delete
			rejectionBody := azuremonitorworkspaces.PrivateEndpointConnection{
				Properties: &azuremonitorworkspaces.PrivateEndpointConnectionProperties{
					PrivateLinkServiceConnectionState: azuremonitorworkspaces.PrivateLinkServiceConnectionState{
						Status:          pointer.To(azuremonitorworkspaces.PrivateEndpointServiceConnectionStatusRejected),
						Description:     pointer.To("Rejected via Terraform destroy"),
						ActionsRequired: pointer.To("None"),
					},
				},
			}

			connectionPath := fmt.Sprintf("%s/privateEndpointConnections/%s", workspaceId.ID(), model.PrivateEndpointConnectionName)

			opts := client.RequestOptions{
				ContentType: "application/json; charset=utf-8",
				ExpectedStatusCodes: []int{
					http.StatusOK,
					http.StatusAccepted,
					http.StatusNoContent,
					http.StatusNotFound,
				},
				HttpMethod: http.MethodPut,
				Path:       connectionPath,
			}

			req, err := workspaceClient.Client.NewRequest(ctx, opts)
			if err != nil {
				return fmt.Errorf("building request: %+v", err)
			}

			if err := req.Marshal(&rejectionBody); err != nil {
				return fmt.Errorf("marshaling request body: %+v", err)
			}

			resp, err := workspaceClient.Client.Execute(ctx, req)
			if err != nil {
				// If the connection is already gone, that's fine
				if resp != nil && resp.StatusCode != http.StatusNotFound {
					metadata.Logger.Infof("rejecting private endpoint connection %q on %s: %+v", model.PrivateEndpointConnectionName, *workspaceId, err)
				}
			}

			return nil
		},
	}
}

// waitForPrivateEndpointConnectionToAppear polls the workspace until the specified private endpoint connection appears.
// This handles eventual consistency when a managed private endpoint is created against this workspace.
func waitForPrivateEndpointConnectionToAppear(ctx context.Context, client *azuremonitorworkspaces.AzureMonitorWorkspacesClient, workspaceId azuremonitorworkspaces.AccountId, connectionName string) (*azuremonitorworkspaces.PrivateEndpointConnection, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(10 * time.Minute)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"pending"},
		Target:                    []string{"found"},
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 2, // Must find the connection twice to confirm it's stable
		Refresh:                   privateEndpointConnectionAppearanceRefreshFunc(ctx, client, workspaceId, connectionName),
		Timeout:                   time.Until(deadline),
	}

	raw, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for private endpoint connection %q to appear: %+v", connectionName, err)
	}

	conn, ok := raw.(*azuremonitorworkspaces.PrivateEndpointConnection)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T returned from state refresh", raw)
	}

	return conn, nil
}

func privateEndpointConnectionAppearanceRefreshFunc(ctx context.Context, client *azuremonitorworkspaces.AzureMonitorWorkspacesClient, workspaceId azuremonitorworkspaces.AccountId, connectionName string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, workspaceId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil, "pending", nil
			}
			return nil, "failed", fmt.Errorf("retrieving %s: %+v", workspaceId, err)
		}

		if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.PrivateEndpointConnections == nil {
			return nil, "pending", nil
		}

		// Look for the connection by name
		for _, conn := range *resp.Model.Properties.PrivateEndpointConnections {
			if conn.Name != nil && *conn.Name == connectionName {
				return &conn, "found", nil
			}
		}

		// Connection not found yet
		return nil, "pending", nil
	}
}
