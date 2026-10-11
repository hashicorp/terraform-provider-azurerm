package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityPostureReference struct {
	ExcludeExtensions *[]string `json:"excludeExtensions,omitempty"`
	Id                string    `json:"id"`
	IsOverridable     *bool     `json:"isOverridable,omitempty"`
}
