package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExternalUserInfo struct {
	EmailId  *string   `json:"emailId,omitempty"`
	FullName *string   `json:"fullName,omitempty"`
	Password *string   `json:"password,omitempty"`
	Roles    *[]string `json:"roles,omitempty"`
	UserName *string   `json:"userName,omitempty"`
}
