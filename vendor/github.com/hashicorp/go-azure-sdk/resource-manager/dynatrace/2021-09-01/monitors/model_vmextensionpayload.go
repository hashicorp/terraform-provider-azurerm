package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMExtensionPayload struct {
	EnvironmentId *string `json:"environmentId,omitempty"`
	IngestionKey  *string `json:"ingestionKey,omitempty"`
}
