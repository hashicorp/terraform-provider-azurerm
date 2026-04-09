package autoscalesettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoscaleSetting struct {
	Enabled                   *bool                      `json:"enabled,omitempty"`
	Name                      *string                    `json:"name,omitempty"`
	Notifications             *[]AutoscaleNotification   `json:"notifications,omitempty"`
	PredictiveAutoscalePolicy *PredictiveAutoscalePolicy `json:"predictiveAutoscalePolicy,omitempty"`
	Profiles                  []AutoscaleProfile         `json:"profiles"`
	TargetResourceLocation    *string                    `json:"targetResourceLocation,omitempty"`
	TargetResourceUri         *string                    `json:"targetResourceUri,omitempty"`
}
