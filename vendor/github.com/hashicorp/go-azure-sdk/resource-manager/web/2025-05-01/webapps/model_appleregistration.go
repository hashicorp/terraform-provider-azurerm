package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppleRegistration struct {
	ClientId                *string `json:"clientId,omitempty"`
	ClientSecretSettingName *string `json:"clientSecretSettingName,omitempty"`
}
