package msixpackage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MsixPackageApplications struct {
	AppId          *string `json:"appId,omitempty"`
	AppUserModelID *string `json:"appUserModelID,omitempty"`
	Description    *string `json:"description,omitempty"`
	FriendlyName   *string `json:"friendlyName,omitempty"`
	IconImageName  *string `json:"iconImageName,omitempty"`
	RawIcon        *string `json:"rawIcon,omitempty"`
	RawPng         *string `json:"rawPng,omitempty"`
}
