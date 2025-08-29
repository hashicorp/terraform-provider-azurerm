package gatewayapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiLicenseInformation struct {
	Name *string `json:"name,omitempty"`
	Url  *string `json:"url,omitempty"`
}
