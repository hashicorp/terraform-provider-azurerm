package notificationhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GcmCredentialProperties struct {
	GcmEndpoint  *string `json:"gcmEndpoint,omitempty"`
	GoogleApiKey *string `json:"googleApiKey,omitempty"`
}
