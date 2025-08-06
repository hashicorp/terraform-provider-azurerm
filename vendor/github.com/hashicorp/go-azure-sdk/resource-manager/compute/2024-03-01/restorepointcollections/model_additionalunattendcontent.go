package restorepointcollections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdditionalUnattendContent struct {
	ComponentName *ComponentNames `json:"componentName,omitempty"`
	Content       *string         `json:"content,omitempty"`
	PassName      *PassNames      `json:"passName,omitempty"`
	SettingName   *SettingNames   `json:"settingName,omitempty"`
}
