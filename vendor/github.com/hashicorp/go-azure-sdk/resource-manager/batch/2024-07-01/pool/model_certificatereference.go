package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateReference struct {
	Id            string                    `json:"id"`
	StoreLocation *CertificateStoreLocation `json:"storeLocation,omitempty"`
	StoreName     *string                   `json:"storeName,omitempty"`
	Visibility    *[]CertificateVisibility  `json:"visibility,omitempty"`
}
