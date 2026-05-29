package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LineChannelProperties struct {
	CallbackURL       *string            `json:"callbackUrl,omitempty"`
	IsValidated       *bool              `json:"isValidated,omitempty"`
	LineRegistrations []LineRegistration `json:"lineRegistrations"`
}
