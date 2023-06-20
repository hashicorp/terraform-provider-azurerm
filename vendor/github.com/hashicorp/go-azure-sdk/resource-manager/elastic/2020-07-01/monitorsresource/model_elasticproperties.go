package monitorsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticProperties struct {
	ElasticCloudDeployment *ElasticCloudDeployment `json:"elasticCloudDeployment,omitempty"`
	ElasticCloudUser       *ElasticCloudUser       `json:"elasticCloudUser,omitempty"`
}
