package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LanguageExtension struct {
	LanguageExtensionImageName *LanguageExtensionImageName `json:"languageExtensionImageName,omitempty"`
	LanguageExtensionName      *LanguageExtensionName      `json:"languageExtensionName,omitempty"`
}
