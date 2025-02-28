package notification

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecipientsContractProperties struct {
	Emails *[]string `json:"emails,omitempty"`
	Users  *[]string `json:"users,omitempty"`
}
