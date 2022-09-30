package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerPropertiesInstanceView struct {
	CurrentState  *ContainerState `json:"currentState,omitempty"`
	Events        *[]Event        `json:"events,omitempty"`
	PreviousState *ContainerState `json:"previousState,omitempty"`
	RestartCount  *int64          `json:"restartCount,omitempty"`
}
