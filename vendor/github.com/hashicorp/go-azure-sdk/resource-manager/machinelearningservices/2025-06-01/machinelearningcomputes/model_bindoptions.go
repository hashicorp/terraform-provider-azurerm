package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BindOptions struct {
	CreateHostPath *bool   `json:"createHostPath,omitempty"`
	Propagation    *string `json:"propagation,omitempty"`
	Selinux        *string `json:"selinux,omitempty"`
}
