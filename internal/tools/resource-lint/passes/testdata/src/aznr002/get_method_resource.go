package aznr002

import (
	"context"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Test Case 3: Using Get() method to retrieve property values
type GetMethodResourceModel struct {
	Name      string `tfschema:"name"`
	AccessKey string `tfschema:"access_key"`
	SecretKey string `tfschema:"secret_key"`
}

type GetMethodResource struct{}

var _ sdk.ResourceWithUpdate = GetMethodResource{}

func (r GetMethodResource) ResourceType() string {
	return "azurerm_get_method_resource"
}

func (r GetMethodResource) ModelObject() interface{} {
	return &GetMethodResourceModel{}
}

func (r GetMethodResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"access_key": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"secret_key": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
}

func (r GetMethodResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r GetMethodResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r GetMethodResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r GetMethodResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r GetMethodResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// Using Get() to access property values - should be detected
			accessKey := metadata.ResourceData.Get("access_key").(string)
			secretKey := metadata.ResourceData.Get("secret_key").(string)

			// Do something with the keys
			_ = accessKey
			_ = secretKey

			return nil
		},
	}
}
