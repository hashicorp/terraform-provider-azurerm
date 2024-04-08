package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DirectLineSpeechChannelProperties struct {
	CognitiveServiceRegion          *string `json:"cognitiveServiceRegion,omitempty"`
	CognitiveServiceResourceId      *string `json:"cognitiveServiceResourceId,omitempty"`
	CognitiveServiceSubscriptionKey *string `json:"cognitiveServiceSubscriptionKey,omitempty"`
	CustomSpeechModelId             *string `json:"customSpeechModelId,omitempty"`
	CustomVoiceDeploymentId         *string `json:"customVoiceDeploymentId,omitempty"`
	IsDefaultBotForCogSvcAccount    *bool   `json:"isDefaultBotForCogSvcAccount,omitempty"`
	IsEnabled                       *bool   `json:"isEnabled,omitempty"`
}
