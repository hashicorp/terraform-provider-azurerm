package environmentversion

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildContext struct {
	ContextUri     string  `json:"contextUri"`
	DockerfilePath *string `json:"dockerfilePath,omitempty"`
}
