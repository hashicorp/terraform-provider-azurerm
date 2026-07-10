package gen

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

func TestTmp_RenderValidateAndDescription(t *testing.T) {
	res := &ir.ResourceIR{
		ModelStructName: "ExampleModel",
		TopLevel: []*ir.Property{
			{TFName: "name", TFType: "TypeString", Required: true, ForceNew: true, Description: "The name of the resource."},
			{TFName: "sku_name", TFType: "TypeString", Optional: true, Description: "The SKU \"tier\"."},
			{TFName: "capacity", TFType: "TypeInt", Optional: true, Description: "Node count."},
			{TFName: "ratio", TFType: "TypeFloat", Optional: true},
			{TFName: "enabled", TFType: "TypeBool", Optional: true, Description: "Whether enabled."},
			{TFName: "kind", TFType: "TypeString", Optional: true, IsEnum: true, EnumValues: []string{"A", "B"}, Description: "The kind."},
			{TFName: "zones", TFType: "TypeList", GoType: "[]string", Optional: true},
			{TFName: "modes", TFType: "TypeList", GoType: "[]string", Optional: true, IsEnum: true, EnumValues: []string{"X", "Y"}},
			{TFName: "ports", TFType: "TypeList", GoType: "[]int64", Optional: true},
			{TFName: "provisioning_state", TFType: "TypeString", Computed: true, Description: "Read-only state."},
			{TFName: "profile", TFType: "TypeList", MaxItems: 1, IsBlock: true, BlockName: "Profile", GoType: "[]Profile", Optional: true, Description: "A profile block."},
		},
		Blocks: []*ir.BlockModel{
			{Name: "Profile", Properties: []*ir.Property{
				{TFName: "vm_size", TFType: "TypeString", Required: true, Description: "The VM size."},
				{TFName: "count", TFType: "TypeInt", Optional: true},
			}},
		},
	}

	fmt.Println("=== ARGUMENTS ===")
	fmt.Println(RenderArguments(res))
	fmt.Println("=== ATTRIBUTES ===")
	fmt.Println(RenderAttributes(res))
}
