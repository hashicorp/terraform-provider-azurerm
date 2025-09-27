package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommandPostBody struct {
	Arguments          *map[string]string `json:"arguments,omitempty"`
	CassandraStopStart *bool              `json:"cassandra-stop-start,omitempty"`
	Command            string             `json:"command"`
	Host               string             `json:"host"`
	Readwrite          *bool              `json:"readwrite,omitempty"`
}
