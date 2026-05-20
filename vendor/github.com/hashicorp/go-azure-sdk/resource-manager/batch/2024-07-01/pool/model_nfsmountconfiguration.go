package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NFSMountConfiguration struct {
	MountOptions      *string `json:"mountOptions,omitempty"`
	RelativeMountPath string  `json:"relativeMountPath"`
	Source            string  `json:"source"`
}
