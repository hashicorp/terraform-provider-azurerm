package desktop

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DesktopProperties struct {
	Description  *string `json:"description,omitempty"`
	FriendlyName *string `json:"friendlyName,omitempty"`
	IconContent  *string `json:"iconContent,omitempty"`
	IconHash     *string `json:"iconHash,omitempty"`
	ObjectId     *string `json:"objectId,omitempty"`
}
