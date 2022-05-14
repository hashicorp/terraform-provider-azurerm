package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func schemaCdnFrontdoorOperator() *pluginsdk.Schema {
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

func schemaCdnFrontdoorOperatorEqualOnly() *pluginsdk.Schema {
	// TODO: if there's only one possible value, and it's defaulted - we don't need to expose this field for now?
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Default:  string(cdn.OperatorEqual),
		ValidateFunc: validation.StringInSlice([]string{
			string(cdn.OperatorEqual),
		}, false),
	}
}

func schemaCdnFrontdoorOperatorRemoteAddress() *pluginsdk.Schema {
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

func schemaCdnFrontdoorOperatorSocketAddress() *pluginsdk.Schema {
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

func schemaCdnFrontdoorNegateCondition() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func schemaCdnFrontdoorMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 25,

		// In some cases it is valid for this to be an empty string
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func schemaCdnFrontdoorServerPortMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: should this be a set?
		Type:     pluginsdk.TypeSet,
		Required: true,
		MaxItems: 2,

		// In some cases it is valid for this to be an empty string
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				"80",
				"443",
			}, false),
		},
	}
}

func schemaCdnFrontdoorSslProtocolMatchValues() *pluginsdk.Schema {
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

func schemaCdnFrontdoorUrlPathConditionMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 25,

		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validate.CdnFrontdoorUrlPathConditionMatchValue,
		},
	}
}

func schemaCdnFrontdoorMatchValuesRequired() *pluginsdk.Schema {
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

func schemaCdnFrontdoorRequestMethodMatchValues() *pluginsdk.Schema {
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

func schemaCdnFrontdoorProtocolMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: if this is MaxItems: 1, should this be a string?
		// WS: This is a list because the up level interface is a list
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,

		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
			Default: "HTTP",
			ValidateFunc: validation.StringInSlice([]string{
				// TODO: are there constants for these?
				// TODO: other APIs use `Http` and `Https`, is that casing consistent in the API?
				"HTTP",
				"HTTPS",
			}, false),
		},
	}
}

func schemaCdnFrontdoorIsDeviceMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,

		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{
				// TODO: are there constants for these?
				// WS: No, these values do not have constants.
				"Mobile",
				"Desktop",
			}, false),
		},
	}
}

func schemaCdnFrontdoorHttpVersionMatchValues() *pluginsdk.Schema {
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

func schemaCdnFrontdoorRuleTransforms() *pluginsdk.Schema {
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
