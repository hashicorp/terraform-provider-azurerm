// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2024-02-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// SslProtocol enumerates the values for ssl protocol.
type SslProtocol string

const (
	// SslProtocolTLSv1 ...
	SslProtocolTLSv1 SslProtocol = "TLSv1"
	// SslProtocolTLSv11 ...
	SslProtocolTLSv11 SslProtocol = "TLSv1.1"
	// SslProtocolTLSv12 ...
	SslProtocolTLSv12 SslProtocol = "TLSv1.2"
)

// PossibleSslProtocolValues returns an array of possible values for the SslProtocol const type.
func PossibleSslProtocolValues() []string {
	return []string{
		string(SslProtocolTLSv1),
		string(SslProtocolTLSv11),
		string(SslProtocolTLSv12),
	}
}

func schemaCdnFrontDoorOperator() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice(waf.PossibleValuesForOperator(),
			false),
	}
}

func schemaCdnFrontDoorOperatorEqualOnly() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(waf.OperatorEqual),
		ValidateFunc: validation.StringInSlice([]string{
			string(waf.OperatorEqual),
		}, false),
	}
}

func schemaCdnFrontDoorOperatorRemoteAddress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(waf.OperatorIPMatch),
		ValidateFunc: validation.StringInSlice([]string{
			string(waf.OperatorAny),
			string(waf.OperatorIPMatch),
			string(waf.OperatorGeoMatch),
		}, false),
	}
}

func schemaCdnFrontDoorOperatorSocketAddress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(waf.OperatorIPMatch),
		ValidateFunc: validation.StringInSlice([]string{
			string(waf.OperatorAny),
			string(waf.OperatorIPMatch),
		}, false),
	}
}

func schemaCdnFrontDoorNegateCondition() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func schemaCdnFrontDoorMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 25,

		Elem: &pluginsdk.Schema{
			// In some cases it is valid for this to be an empty string
			Type: pluginsdk.TypeString,
		},
	}
}

func schemaCdnFrontDoorServerPortMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Required: true,
		MaxItems: 2,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				"80",
				"443",
			}, false),
		},
	}
}

func schemaCdnFrontDoorSslProtocolMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Required: true,
		MaxItems: 3,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(PossibleSslProtocolValues(),
				false),
		},
	}
}

func schemaCdnFrontDoorUrlPathConditionMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 25,

		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validate.CdnFrontDoorUrlPathConditionMatchValue,
		},
	}
}

func schemaCdnFrontDoorMatchValuesRequired() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 25,

		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func schemaCdnFrontDoorRequestMethodMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Required: true,
		MaxItems: 7,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				"GET",
				"POST",
				"PUT",
				"DELETE",
				"HEAD",
				"OPTIONS",
				"TRACE",
			}, false),
		},
	}
}

func schemaCdnFrontDoorProtocolMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,

		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
			Default: "HTTP",
			ValidateFunc: validation.StringInSlice([]string{
				"HTTP",
				"HTTPS",
			}, false),
		},
	}
}

func schemaCdnFrontDoorIsDeviceMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				"Mobile",
				"Desktop",
			}, false),
		},
	}
}

func schemaCdnFrontDoorHttpVersionMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Required: true,
		MaxItems: 4,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				"2.0",
				"1.1",
				"1.0",
				"0.9",
			}, false),
		},
	}
}

func schemaCdnFrontDoorRuleTransforms() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		MaxItems: 4,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(waf.PossibleValuesForTransformType(),
				false),
		},
	}
}
