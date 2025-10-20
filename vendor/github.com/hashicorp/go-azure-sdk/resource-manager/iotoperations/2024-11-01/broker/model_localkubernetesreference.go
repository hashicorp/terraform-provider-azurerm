package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalKubernetesReference struct {
	ApiGroup *string `json:"apiGroup,omitempty"`
	Kind     string  `json:"kind"`
	Name     string  `json:"name"`
}
