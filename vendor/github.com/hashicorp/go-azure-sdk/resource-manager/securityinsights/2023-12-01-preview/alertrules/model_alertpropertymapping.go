package alertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertPropertyMapping struct {
	AlertProperty *AlertProperty `json:"alertProperty,omitempty"`
	Value         *string        `json:"value,omitempty"`
}
