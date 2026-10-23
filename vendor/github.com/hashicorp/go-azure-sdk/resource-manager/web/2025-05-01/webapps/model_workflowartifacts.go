package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowArtifacts struct {
	AppSettings   *interface{}            `json:"appSettings,omitempty"`
	Files         *map[string]interface{} `json:"files,omitempty"`
	FilesToDelete *[]string               `json:"filesToDelete,omitempty"`
}
