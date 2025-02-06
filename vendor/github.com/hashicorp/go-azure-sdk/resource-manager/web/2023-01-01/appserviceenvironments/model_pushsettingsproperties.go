package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PushSettingsProperties struct {
	DynamicTagsJson   *string `json:"dynamicTagsJson,omitempty"`
	IsPushEnabled     bool    `json:"isPushEnabled"`
	TagWhitelistJson  *string `json:"tagWhitelistJson,omitempty"`
	TagsRequiringAuth *string `json:"tagsRequiringAuth,omitempty"`
}
