package application

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationProperties struct {
	ApplicationType          *RemoteApplicationType `json:"applicationType,omitempty"`
	CommandLineArguments     *string                `json:"commandLineArguments,omitempty"`
	CommandLineSetting       CommandLineSetting     `json:"commandLineSetting"`
	Description              *string                `json:"description,omitempty"`
	FilePath                 *string                `json:"filePath,omitempty"`
	FriendlyName             *string                `json:"friendlyName,omitempty"`
	IconContent              *string                `json:"iconContent,omitempty"`
	IconHash                 *string                `json:"iconHash,omitempty"`
	IconIndex                *int64                 `json:"iconIndex,omitempty"`
	IconPath                 *string                `json:"iconPath,omitempty"`
	MsixPackageApplicationId *string                `json:"msixPackageApplicationId,omitempty"`
	MsixPackageFamilyName    *string                `json:"msixPackageFamilyName,omitempty"`
	ObjectId                 *string                `json:"objectId,omitempty"`
	ShowInPortal             *bool                  `json:"showInPortal,omitempty"`
}
