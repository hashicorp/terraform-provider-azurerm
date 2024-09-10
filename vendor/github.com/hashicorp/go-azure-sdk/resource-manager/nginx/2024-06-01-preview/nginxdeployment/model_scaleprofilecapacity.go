package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleProfileCapacity struct {
	Max int64 `json:"max"`
	Min int64 `json:"min"`
}
