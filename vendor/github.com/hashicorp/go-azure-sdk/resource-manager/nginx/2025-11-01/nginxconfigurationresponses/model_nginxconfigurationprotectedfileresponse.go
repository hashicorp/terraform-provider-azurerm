package nginxconfigurationresponses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxConfigurationProtectedFileResponse struct {
	ContentHash *string `json:"contentHash,omitempty"`
	VirtualPath *string `json:"virtualPath,omitempty"`
}
