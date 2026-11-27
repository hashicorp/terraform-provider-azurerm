package msiximage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MsixPackageDependencies struct {
	DependencyName *string `json:"dependencyName,omitempty"`
	MinVersion     *string `json:"minVersion,omitempty"`
	Publisher      *string `json:"publisher,omitempty"`
}
