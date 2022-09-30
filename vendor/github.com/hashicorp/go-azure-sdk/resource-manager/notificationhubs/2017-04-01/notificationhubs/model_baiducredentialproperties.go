package notificationhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaiduCredentialProperties struct {
	BaiduApiKey    *string `json:"baiduApiKey,omitempty"`
	BaiduEndPoint  *string `json:"baiduEndPoint,omitempty"`
	BaiduSecretKey *string `json:"baiduSecretKey,omitempty"`
}
