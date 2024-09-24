package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Organization struct {
	Parallelism *int64    `json:"parallelism,omitempty"`
	Projects    *[]string `json:"projects,omitempty"`
	Url         string    `json:"url"`
}
