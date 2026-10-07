package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteConfig struct {
	AcrUseManagedIdentityCreds             *bool                             `json:"acrUseManagedIdentityCreds,omitempty"`
	AcrUserManagedIdentityID               *string                           `json:"acrUserManagedIdentityID,omitempty"`
	AlwaysOn                               *bool                             `json:"alwaysOn,omitempty"`
	ApiDefinition                          *ApiDefinitionInfo                `json:"apiDefinition,omitempty"`
	ApiManagementConfig                    *ApiManagementConfig              `json:"apiManagementConfig,omitempty"`
	AppCommandLine                         *string                           `json:"appCommandLine,omitempty"`
	AppSettings                            *[]NameValuePair                  `json:"appSettings,omitempty"`
	AutoHealEnabled                        *bool                             `json:"autoHealEnabled,omitempty"`
	AutoHealRules                          *AutoHealRules                    `json:"autoHealRules,omitempty"`
	AutoSwapSlotName                       *string                           `json:"autoSwapSlotName,omitempty"`
	AzureStorageAccounts                   *map[string]AzureStorageInfoValue `json:"azureStorageAccounts,omitempty"`
	ConnectionStrings                      *[]ConnStringInfo                 `json:"connectionStrings,omitempty"`
	Cors                                   *CorsSettings                     `json:"cors,omitempty"`
	DefaultDocuments                       *[]string                         `json:"defaultDocuments,omitempty"`
	DetailedErrorLoggingEnabled            *bool                             `json:"detailedErrorLoggingEnabled,omitempty"`
	DocumentRoot                           *string                           `json:"documentRoot,omitempty"`
	ElasticWebAppScaleLimit                *int64                            `json:"elasticWebAppScaleLimit,omitempty"`
	Experiments                            *Experiments                      `json:"experiments,omitempty"`
	FtpsState                              *FtpsState                        `json:"ftpsState,omitempty"`
	FunctionAppScaleLimit                  *int64                            `json:"functionAppScaleLimit,omitempty"`
	FunctionsRuntimeScaleMonitoringEnabled *bool                             `json:"functionsRuntimeScaleMonitoringEnabled,omitempty"`
	HTTP20Enabled                          *bool                             `json:"http20Enabled,omitempty"`
	HTTPLoggingEnabled                     *bool                             `json:"httpLoggingEnabled,omitempty"`
	HandlerMappings                        *[]HandlerMapping                 `json:"handlerMappings,omitempty"`
	HealthCheckPath                        *string                           `json:"healthCheckPath,omitempty"`
	IPSecurityRestrictions                 *[]IPSecurityRestriction          `json:"ipSecurityRestrictions,omitempty"`
	IPSecurityRestrictionsDefaultAction    *DefaultAction                    `json:"ipSecurityRestrictionsDefaultAction,omitempty"`
	JavaContainer                          *string                           `json:"javaContainer,omitempty"`
	JavaContainerVersion                   *string                           `json:"javaContainerVersion,omitempty"`
	JavaVersion                            *string                           `json:"javaVersion,omitempty"`
	KeyVaultReferenceIdentity              *string                           `json:"keyVaultReferenceIdentity,omitempty"`
	Limits                                 *SiteLimits                       `json:"limits,omitempty"`
	LinuxFxVersion                         *string                           `json:"linuxFxVersion,omitempty"`
	LoadBalancing                          *SiteLoadBalancing                `json:"loadBalancing,omitempty"`
	LocalMySqlEnabled                      *bool                             `json:"localMySqlEnabled,omitempty"`
	LogsDirectorySizeLimit                 *int64                            `json:"logsDirectorySizeLimit,omitempty"`
	MachineKey                             *SiteMachineKey                   `json:"machineKey,omitempty"`
	ManagedPipelineMode                    *ManagedPipelineMode              `json:"managedPipelineMode,omitempty"`
	ManagedServiceIdentityId               *int64                            `json:"managedServiceIdentityId,omitempty"`
	Metadata                               *[]NameValuePair                  `json:"metadata,omitempty"`
	MinTlsCipherSuite                      *TlsCipherSuites                  `json:"minTlsCipherSuite,omitempty"`
	MinTlsVersion                          *SupportedTlsVersions             `json:"minTlsVersion,omitempty"`
	MinimumElasticInstanceCount            *int64                            `json:"minimumElasticInstanceCount,omitempty"`
	NetFrameworkVersion                    *string                           `json:"netFrameworkVersion,omitempty"`
	NodeVersion                            *string                           `json:"nodeVersion,omitempty"`
	NumberOfWorkers                        *int64                            `json:"numberOfWorkers,omitempty"`
	PhpVersion                             *string                           `json:"phpVersion,omitempty"`
	PowerShellVersion                      *string                           `json:"powerShellVersion,omitempty"`
	PreWarmedInstanceCount                 *int64                            `json:"preWarmedInstanceCount,omitempty"`
	PublicNetworkAccess                    *string                           `json:"publicNetworkAccess,omitempty"`
	PublishingUsername                     *string                           `json:"publishingUsername,omitempty"`
	Push                                   *PushSettings                     `json:"push,omitempty"`
	PythonVersion                          *string                           `json:"pythonVersion,omitempty"`
	RemoteDebuggingEnabled                 *bool                             `json:"remoteDebuggingEnabled,omitempty"`
	RemoteDebuggingVersion                 *string                           `json:"remoteDebuggingVersion,omitempty"`
	RequestTracingEnabled                  *bool                             `json:"requestTracingEnabled,omitempty"`
	RequestTracingExpirationTime           *string                           `json:"requestTracingExpirationTime,omitempty"`
	ScmIPSecurityRestrictions              *[]IPSecurityRestriction          `json:"scmIpSecurityRestrictions,omitempty"`
	ScmIPSecurityRestrictionsDefaultAction *DefaultAction                    `json:"scmIpSecurityRestrictionsDefaultAction,omitempty"`
	ScmIPSecurityRestrictionsUseMain       *bool                             `json:"scmIpSecurityRestrictionsUseMain,omitempty"`
	ScmMinTlsVersion                       *SupportedTlsVersions             `json:"scmMinTlsVersion,omitempty"`
	ScmType                                *ScmType                          `json:"scmType,omitempty"`
	TracingOptions                         *string                           `json:"tracingOptions,omitempty"`
	Use32BitWorkerProcess                  *bool                             `json:"use32BitWorkerProcess,omitempty"`
	VirtualApplications                    *[]VirtualApplication             `json:"virtualApplications,omitempty"`
	VnetName                               *string                           `json:"vnetName,omitempty"`
	VnetPrivatePortsCount                  *int64                            `json:"vnetPrivatePortsCount,omitempty"`
	VnetRouteAllEnabled                    *bool                             `json:"vnetRouteAllEnabled,omitempty"`
	WebSocketsEnabled                      *bool                             `json:"webSocketsEnabled,omitempty"`
	WebsiteTimeZone                        *string                           `json:"websiteTimeZone,omitempty"`
	WindowsFxVersion                       *string                           `json:"windowsFxVersion,omitempty"`
	XManagedServiceIdentityId              *int64                            `json:"xManagedServiceIdentityId,omitempty"`
}

func (o *SiteConfig) GetRequestTracingExpirationTimeAsTime() (*time.Time, error) {
	if o.RequestTracingExpirationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RequestTracingExpirationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SiteConfig) SetRequestTracingExpirationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RequestTracingExpirationTime = &formatted
}
