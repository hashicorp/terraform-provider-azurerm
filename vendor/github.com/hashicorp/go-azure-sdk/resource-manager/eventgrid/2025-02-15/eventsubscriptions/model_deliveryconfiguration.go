package eventsubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeliveryConfiguration struct {
	DeliveryMode *DeliveryMode `json:"deliveryMode,omitempty"`
	Push         *PushInfo     `json:"push,omitempty"`
	Queue        *QueueInfo    `json:"queue,omitempty"`
}
