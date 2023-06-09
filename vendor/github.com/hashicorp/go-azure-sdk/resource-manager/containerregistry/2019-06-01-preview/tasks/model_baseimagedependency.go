package tasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseImageDependency struct {
	Digest     *string                  `json:"digest,omitempty"`
	Registry   *string                  `json:"registry,omitempty"`
	Repository *string                  `json:"repository,omitempty"`
	Tag        *string                  `json:"tag,omitempty"`
	Type       *BaseImageDependencyType `json:"type,omitempty"`
}
