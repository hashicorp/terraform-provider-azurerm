package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueueScaleRule struct {
	Auth        *[]ScaleRuleAuth `json:"auth,omitempty"`
	QueueLength *int64           `json:"queueLength,omitempty"`
	QueueName   *string          `json:"queueName,omitempty"`
}
