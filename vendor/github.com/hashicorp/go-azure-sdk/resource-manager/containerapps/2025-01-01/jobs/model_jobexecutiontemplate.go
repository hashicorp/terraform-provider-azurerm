package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobExecutionTemplate struct {
	Containers     *[]JobExecutionContainer `json:"containers,omitempty"`
	InitContainers *[]JobExecutionContainer `json:"initContainers,omitempty"`
}
