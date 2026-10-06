package hubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationHubProperties struct {
	AdmCredential         *AdmCredential                             `json:"admCredential,omitempty"`
	ApnsCredential        *ApnsCredential                            `json:"apnsCredential,omitempty"`
	AuthorizationRules    *[]SharedAccessAuthorizationRuleProperties `json:"authorizationRules,omitempty"`
	BaiduCredential       *BaiduCredential                           `json:"baiduCredential,omitempty"`
	BrowserCredential     *BrowserCredential                         `json:"browserCredential,omitempty"`
	DailyMaxActiveDevices *int64                                     `json:"dailyMaxActiveDevices,omitempty"`
	GcmCredential         *GcmCredential                             `json:"gcmCredential,omitempty"`
	MpnsCredential        *MpnsCredential                            `json:"mpnsCredential,omitempty"`
	Name                  *string                                    `json:"name,omitempty"`
	RegistrationTtl       *string                                    `json:"registrationTtl,omitempty"`
	WnsCredential         *WnsCredential                             `json:"wnsCredential,omitempty"`
	XiaomiCredential      *XiaomiCredential                          `json:"xiaomiCredential,omitempty"`
}
