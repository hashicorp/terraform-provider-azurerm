package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobConfiguration struct {
	EventTriggerConfig    *JobConfigurationEventTriggerConfig    `json:"eventTriggerConfig,omitempty"`
	IdentitySettings      *[]IdentitySettings                    `json:"identitySettings,omitempty"`
	ManualTriggerConfig   *JobConfigurationManualTriggerConfig   `json:"manualTriggerConfig,omitempty"`
	Registries            *[]RegistryCredentials                 `json:"registries,omitempty"`
	ReplicaRetryLimit     *int64                                 `json:"replicaRetryLimit,omitempty"`
	ReplicaTimeout        int64                                  `json:"replicaTimeout"`
	ScheduleTriggerConfig *JobConfigurationScheduleTriggerConfig `json:"scheduleTriggerConfig,omitempty"`
	Secrets               *[]Secret                              `json:"secrets,omitempty"`
	TriggerType           TriggerType                            `json:"triggerType"`
}
