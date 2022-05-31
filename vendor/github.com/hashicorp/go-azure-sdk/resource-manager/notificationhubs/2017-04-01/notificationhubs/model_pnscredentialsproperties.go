package notificationhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PnsCredentialsProperties struct {
	AdmCredential   *AdmCredential   `json:"admCredential,omitempty"`
	ApnsCredential  *ApnsCredential  `json:"apnsCredential,omitempty"`
	BaiduCredential *BaiduCredential `json:"baiduCredential,omitempty"`
	GcmCredential   *GcmCredential   `json:"gcmCredential,omitempty"`
	MpnsCredential  *MpnsCredential  `json:"mpnsCredential,omitempty"`
	WnsCredential   *WnsCredential   `json:"wnsCredential,omitempty"`
}
