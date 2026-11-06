package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeepCreatedOrigin struct {
	Name       string                       `json:"name"`
	Properties *DeepCreatedOriginProperties `json:"properties,omitempty"`
}
