package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateSecretOutput struct {
	Message *string `json:"message,omitempty"`
	Status  *Status `json:"status,omitempty"`
}
