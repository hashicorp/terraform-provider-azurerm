package azure

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func SchemaApiManagementName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementServiceName,
	}
}

func SchemaApiManagementDataSourceName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validate.ApiManagementServiceName,
	}
}

func SchemaApiManagementProductName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.ApiManagementProductName,
	}
}

func SchemaApiManagementProductDataSourceName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validate.ApiManagementProductName,
	}
}
