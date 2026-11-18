package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserProperties struct {
	PublishingPassword         *string `json:"publishingPassword,omitempty"`
	PublishingPasswordHash     *string `json:"publishingPasswordHash,omitempty"`
	PublishingPasswordHashSalt *string `json:"publishingPasswordHashSalt,omitempty"`
	PublishingUserName         string  `json:"publishingUserName"`
	ScmUri                     *string `json:"scmUri,omitempty"`
}
