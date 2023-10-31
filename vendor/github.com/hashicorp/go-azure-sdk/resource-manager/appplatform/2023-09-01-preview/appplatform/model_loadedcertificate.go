package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadedCertificate struct {
	LoadTrustStore *bool  `json:"loadTrustStore,omitempty"`
	ResourceId     string `json:"resourceId"`
}
