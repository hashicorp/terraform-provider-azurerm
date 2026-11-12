package goldengateconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KafkaBootstrapServer struct {
	Host string `json:"host"`
	Port *int64 `json:"port,omitempty"`
}
