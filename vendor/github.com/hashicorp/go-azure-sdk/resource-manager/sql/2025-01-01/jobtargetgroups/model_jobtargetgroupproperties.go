package jobtargetgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobTargetGroupProperties struct {
	Members []JobTarget `json:"members"`
}
