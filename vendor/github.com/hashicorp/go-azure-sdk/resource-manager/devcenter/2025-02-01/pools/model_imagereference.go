package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageReference struct {
	ExactVersion *string `json:"exactVersion,omitempty"`
	Id           *string `json:"id,omitempty"`
}
