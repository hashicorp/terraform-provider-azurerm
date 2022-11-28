package notificationhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationHubProperties struct {
	AdmCredential      *AdmCredential                             `json:"admCredential"`
	ApnsCredential     *ApnsCredential                            `json:"apnsCredential"`
	AuthorizationRules *[]SharedAccessAuthorizationRuleProperties `json:"authorizationRules,omitempty"`
	BaiduCredential    *BaiduCredential                           `json:"baiduCredential"`
	GcmCredential      *GcmCredential                             `json:"gcmCredential"`
	MpnsCredential     *MpnsCredential                            `json:"mpnsCredential"`
	Name               *string                                    `json:"name,omitempty"`
	RegistrationTtl    *string                                    `json:"registrationTtl,omitempty"`
	WnsCredential      *WnsCredential                             `json:"wnsCredential"`
}
