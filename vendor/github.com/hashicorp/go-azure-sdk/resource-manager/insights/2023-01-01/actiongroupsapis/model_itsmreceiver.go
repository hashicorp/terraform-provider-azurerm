package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ItsmReceiver struct {
	ConnectionId        string `json:"connectionId"`
	Name                string `json:"name"`
	Region              string `json:"region"`
	TicketConfiguration string `json:"ticketConfiguration"`
	WorkspaceId         string `json:"workspaceId"`
}
