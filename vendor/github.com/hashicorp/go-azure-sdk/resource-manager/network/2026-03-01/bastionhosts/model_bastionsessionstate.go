package bastionhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionSessionState struct {
	Message   *string `json:"message,omitempty"`
	SessionId *string `json:"sessionId,omitempty"`
	State     *string `json:"state,omitempty"`
}
