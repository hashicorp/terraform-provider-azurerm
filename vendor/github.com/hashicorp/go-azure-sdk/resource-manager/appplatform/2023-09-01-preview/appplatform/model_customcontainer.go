package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomContainer struct {
	Args                    *[]string                `json:"args,omitempty"`
	Command                 *[]string                `json:"command,omitempty"`
	ContainerImage          *string                  `json:"containerImage,omitempty"`
	ImageRegistryCredential *ImageRegistryCredential `json:"imageRegistryCredential,omitempty"`
	LanguageFramework       *string                  `json:"languageFramework,omitempty"`
	Server                  *string                  `json:"server,omitempty"`
}
