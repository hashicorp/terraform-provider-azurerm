// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdnfrontdoorruleconditions

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
)

type CdnFrontDoorConditionParameters struct {
	Name       rules.MatchVariable
	TypeName   rules.DeliveryRuleConditionParametersType
	ConfigName string
}

type CdnFrontDoorCondtionsMappings struct {
	ClientPort       CdnFrontDoorConditionParameters
	Cookies          CdnFrontDoorConditionParameters
	HostName         CdnFrontDoorConditionParameters
	HttpVersion      CdnFrontDoorConditionParameters
	IsDevice         CdnFrontDoorConditionParameters
	PostArgs         CdnFrontDoorConditionParameters
	QueryString      CdnFrontDoorConditionParameters
	RemoteAddress    CdnFrontDoorConditionParameters
	RequestBody      CdnFrontDoorConditionParameters
	RequestHeader    CdnFrontDoorConditionParameters
	RequestMethod    CdnFrontDoorConditionParameters
	RequestScheme    CdnFrontDoorConditionParameters
	RequestUri       CdnFrontDoorConditionParameters
	ServerPort       CdnFrontDoorConditionParameters
	SocketAddress    CdnFrontDoorConditionParameters
	SslProtocol      CdnFrontDoorConditionParameters
	UrlFileExtension CdnFrontDoorConditionParameters
	UrlFilename      CdnFrontDoorConditionParameters
	UrlPath          CdnFrontDoorConditionParameters
}

type normalizedSelector struct {
	name  *string
	value *string
}

type normalizedCondition struct {
	selector        *normalizedSelector
	operator        string
	negateCondition *bool
	matchValues     *[]string
	transforms      *[]rules.Transform
}
