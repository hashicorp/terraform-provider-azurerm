package application

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationPatchProperties struct {
	ApplicationType          *RemoteApplicationType `json:"applicationType,omitempty"`
	CommandLineArguments     *string                `json:"commandLineArguments,omitempty"`
	CommandLineSetting       *CommandLineSetting    `json:"commandLineSetting,omitempty"`
	Description              *string                `json:"description,omitempty"`
	FilePath                 *string                `json:"filePath,omitempty"`
	FriendlyName             *string                `json:"friendlyName,omitempty"`
	IconIndex                *int64                 `json:"iconIndex,omitempty"`
	IconPath                 *string                `json:"iconPath,omitempty"`
	MsixPackageApplicationId *string                `json:"msixPackageApplicationId,omitempty"`
	MsixPackageFamilyName    *string                `json:"msixPackageFamilyName,omitempty"`
	ShowInPortal             *bool                  `json:"showInPortal,omitempty"`
}
