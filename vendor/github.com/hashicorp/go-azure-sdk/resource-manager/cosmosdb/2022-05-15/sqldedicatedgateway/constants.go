package sqldedicatedgateway

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceSize string

const (
	ServiceSizeCosmosPointDEights  ServiceSize = "Cosmos.D8s"
	ServiceSizeCosmosPointDFours   ServiceSize = "Cosmos.D4s"
	ServiceSizeCosmosPointDOneSixs ServiceSize = "Cosmos.D16s"
)

func PossibleValuesForServiceSize() []string {
	return []string{
		string(ServiceSizeCosmosPointDEights),
		string(ServiceSizeCosmosPointDFours),
		string(ServiceSizeCosmosPointDOneSixs),
	}
}

func parseServiceSize(input string) (*ServiceSize, error) {
	vals := map[string]ServiceSize{
		"cosmos.d8s":  ServiceSizeCosmosPointDEights,
		"cosmos.d4s":  ServiceSizeCosmosPointDFours,
		"cosmos.d16s": ServiceSizeCosmosPointDOneSixs,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceSize(input)
	return &out, nil
}

type ServiceStatus string

const (
	ServiceStatusCreating ServiceStatus = "Creating"
	ServiceStatusDeleting ServiceStatus = "Deleting"
	ServiceStatusError    ServiceStatus = "Error"
	ServiceStatusRunning  ServiceStatus = "Running"
	ServiceStatusStopped  ServiceStatus = "Stopped"
	ServiceStatusUpdating ServiceStatus = "Updating"
)

func PossibleValuesForServiceStatus() []string {
	return []string{
		string(ServiceStatusCreating),
		string(ServiceStatusDeleting),
		string(ServiceStatusError),
		string(ServiceStatusRunning),
		string(ServiceStatusStopped),
		string(ServiceStatusUpdating),
	}
}

func parseServiceStatus(input string) (*ServiceStatus, error) {
	vals := map[string]ServiceStatus{
		"creating": ServiceStatusCreating,
		"deleting": ServiceStatusDeleting,
		"error":    ServiceStatusError,
		"running":  ServiceStatusRunning,
		"stopped":  ServiceStatusStopped,
		"updating": ServiceStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceStatus(input)
	return &out, nil
}

type ServiceType string

const (
	ServiceTypeDataTransfer             ServiceType = "DataTransfer"
	ServiceTypeGraphAPICompute          ServiceType = "GraphAPICompute"
	ServiceTypeMaterializedViewsBuilder ServiceType = "MaterializedViewsBuilder"
	ServiceTypeSqlDedicatedGateway      ServiceType = "SqlDedicatedGateway"
)

func PossibleValuesForServiceType() []string {
	return []string{
		string(ServiceTypeDataTransfer),
		string(ServiceTypeGraphAPICompute),
		string(ServiceTypeMaterializedViewsBuilder),
		string(ServiceTypeSqlDedicatedGateway),
	}
}

func parseServiceType(input string) (*ServiceType, error) {
	vals := map[string]ServiceType{
		"datatransfer":             ServiceTypeDataTransfer,
		"graphapicompute":          ServiceTypeGraphAPICompute,
		"materializedviewsbuilder": ServiceTypeMaterializedViewsBuilder,
		"sqldedicatedgateway":      ServiceTypeSqlDedicatedGateway,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceType(input)
	return &out, nil
}
