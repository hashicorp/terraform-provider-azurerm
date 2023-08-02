package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPrivateLinkServiceConnectionStateProperty struct {
	ActionsRequired *PrivateLinkServiceConnectionStateActionsRequire `json:"actionsRequired,omitempty"`
	Description     string                                           `json:"description"`
	Status          PrivateLinkServiceConnectionStateStatus          `json:"status"`
}
