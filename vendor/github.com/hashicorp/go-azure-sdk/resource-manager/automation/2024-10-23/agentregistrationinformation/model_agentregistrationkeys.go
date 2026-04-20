package agentregistrationinformation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentRegistrationKeys struct {
	Primary   *string `json:"primary,omitempty"`
	Secondary *string `json:"secondary,omitempty"`
}
