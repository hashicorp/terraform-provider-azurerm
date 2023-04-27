package autoscalesettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleAction struct {
	Cooldown  string         `json:"cooldown"`
	Direction ScaleDirection `json:"direction"`
	Type      ScaleType      `json:"type"`
	Value     *string        `json:"value,omitempty"`
}
