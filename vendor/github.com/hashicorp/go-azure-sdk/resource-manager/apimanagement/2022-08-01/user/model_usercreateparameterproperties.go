package user

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserCreateParameterProperties struct {
	AppType      *AppType                `json:"appType,omitempty"`
	Confirmation *Confirmation           `json:"confirmation,omitempty"`
	Email        string                  `json:"email"`
	FirstName    string                  `json:"firstName"`
	Identities   *[]UserIdentityContract `json:"identities,omitempty"`
	LastName     string                  `json:"lastName"`
	Note         *string                 `json:"note,omitempty"`
	Password     *string                 `json:"password,omitempty"`
	State        *UserState              `json:"state,omitempty"`
}
