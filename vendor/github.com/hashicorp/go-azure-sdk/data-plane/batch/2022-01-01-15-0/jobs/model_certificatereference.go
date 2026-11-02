package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateReference struct {
	StoreLocation       *CertificateStoreLocation `json:"storeLocation,omitempty"`
	StoreName           *string                   `json:"storeName,omitempty"`
	Thumbprint          string                    `json:"thumbprint"`
	ThumbprintAlgorithm string                    `json:"thumbprintAlgorithm"`
	Visibility          *[]CertificateVisibility  `json:"visibility,omitempty"`
}
