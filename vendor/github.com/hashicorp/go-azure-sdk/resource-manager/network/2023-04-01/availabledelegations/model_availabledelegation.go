package availabledelegations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableDelegation struct {
	Actions     *[]string `json:"actions,omitempty"`
	Id          *string   `json:"id,omitempty"`
	Name        *string   `json:"name,omitempty"`
	ServiceName *string   `json:"serviceName,omitempty"`
	Type        *string   `json:"type,omitempty"`
}
