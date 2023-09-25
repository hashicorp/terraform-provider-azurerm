package appliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupportedVersionCatalogVersionData struct {
	Audience *string `json:"audience,omitempty"`
	Catalog  *string `json:"catalog,omitempty"`
	Offer    *string `json:"offer,omitempty"`
	Version  *string `json:"version,omitempty"`
}
