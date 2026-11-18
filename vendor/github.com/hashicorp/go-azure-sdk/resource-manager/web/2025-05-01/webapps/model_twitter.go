package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Twitter struct {
	Enabled      *bool                `json:"enabled,omitempty"`
	Registration *TwitterRegistration `json:"registration,omitempty"`
}
