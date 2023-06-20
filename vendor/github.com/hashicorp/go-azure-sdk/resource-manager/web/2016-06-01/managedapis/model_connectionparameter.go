package managedapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionParameter struct {
	OAuthSettings *ApiOAuthSettings        `json:"oAuthSettings,omitempty"`
	Type          *ConnectionParameterType `json:"type,omitempty"`
}
