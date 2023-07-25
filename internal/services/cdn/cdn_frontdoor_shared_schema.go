// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func schemaCdnFrontDoorOperator() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorAny),
			string(cdn.OperatorEqual),
			string(cdn.OperatorContains),
			string(cdn.OperatorBeginsWith),
			string(cdn.OperatorEndsWith),
			string(cdn.OperatorLessThan),
			string(cdn.OperatorLessThanOrEqual),
			string(cdn.OperatorGreaterThan),
			string(cdn.OperatorGreaterThanOrEqual),
			string(cdn.OperatorRegEx),
		}, false),
	}
}

func schemaCdnFrontDoorOperatorEqualOnly() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(cdn.OperatorEqual),
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorEqual),
		}, false),
	}
}

func schemaCdnFrontDoorOperatorRemoteAddress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(cdn.OperatorIPMatch),
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorAny),
			string(cdn.OperatorIPMatch),
			string(cdn.OperatorGeoMatch),
		}, false),
	}
}

func schemaCdnFrontDoorOperatorSocketAddress() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(cdn.OperatorIPMatch),
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorAny),
			string(cdn.OperatorIPMatch),
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
			ValidateFunc: validation.StringInSlice([]string{
				string(cdn.SslProtocolTLSv1),
				string(cdn.SslProtocolTLSv11),
				string(cdn.SslProtocolTLSv12),
			}, false),
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
			ValidateFunc: validation.StringInSlice([]string{
				string(cdn.TransformLowercase),
				string(cdn.TransformRemoveNulls),
				string(cdn.TransformTrim),
				string(cdn.TransformUppercase),
				string(cdn.TransformURLDecode),
				string(cdn.TransformURLEncode),
			}, false),
		},
	}
}
