package service

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceStatistics struct {
	Counters ServiceCounters `json:"counters"`
	Limits   ServiceLimits   `json:"limits"`
}
