package notificationhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PnsCredentialsProperties struct {
	AdmCredential   *AdmCredential   `json:"admCredential"`
	ApnsCredential  *ApnsCredential  `json:"apnsCredential"`
	BaiduCredential *BaiduCredential `json:"baiduCredential"`
	GcmCredential   *GcmCredential   `json:"gcmCredential"`
	MpnsCredential  *MpnsCredential  `json:"mpnsCredential"`
	WnsCredential   *WnsCredential   `json:"wnsCredential"`
}
