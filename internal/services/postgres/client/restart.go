// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/serverrestart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/servers"
)

// RestartServer restarts a PostgreSQL Flexible Server.
// It checks the server's state and waits for it to be in a state that allows for a restart.
// If the server is already in a restarting state, it waits for the operation to complete.
// If the server is disabled or stopped, it returns an error.
// If the server is in a valid state, it initiates the restart operation.
func (c *Client) RestartServer(ctx context.Context, id serverrestart.FlexibleServerId) error {
	serverID := servers.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)

	for {
		server, err := c.FlexibleServersClient.Get(ctx, serverID)
		if err != nil {
			return fmt.Errorf("getting server %s: %+v", id, err)
		}

		if server.Model == nil || server.Model.Properties == nil || server.Model.Properties.State == nil {
			return fmt.Errorf("server %s has no state", id)
		}

		state := *server.Model.Properties.State
		switch state {
		case servers.ServerStateDisabled, servers.ServerStateStopped, servers.ServerStateDropping:
			return fmt.Errorf("cannot restart server: %s is disabled or stopped", id)
		case servers.ServerStateStarting, servers.ServerStateStopping, servers.ServerStateUpdating, "Restarting":
			// Server is already in a restarting state, wait for it to complete
			// Restarting is a valid state but not defined in the enum: https://github.com/Azure/azure-rest-api-specs/issues/33753
			log.Printf("[DEBUG] Server %s is in state %s, waiting for it to complete", id, state)
			time.Sleep(10 * time.Second)
			continue
		default:
			// Server is ready, we can proceed with the restart
			if err := c.ServerRestartClient.ServersRestartThenPoll(ctx, id, serverrestart.RestartParameter{}); err != nil {
				return fmt.Errorf("restarting server %s: %+v", id, err)
			}
			return nil
		}
	}
}
