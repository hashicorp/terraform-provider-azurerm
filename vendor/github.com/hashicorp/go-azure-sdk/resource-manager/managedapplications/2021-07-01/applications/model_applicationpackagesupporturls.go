package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationPackageSupportUrls struct {
	GovernmentCloud *string `json:"governmentCloud,omitempty"`
	PublicAzure     *string `json:"publicAzure,omitempty"`
}
