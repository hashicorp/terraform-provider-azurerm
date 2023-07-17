package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/generator-typed-resource/logic"
)

func statistic() {
	var typedR, typedDS, untypedR, untypedDS int
	for _, reg := range provider.SupportedTypedServices() {
		typedR += len(reg.Resources())
		typedDS += len(reg.DataSources())
	}

	for _, reg := range provider.SupportedUntypedServices() {
		untypedR += len(reg.SupportedResources())
		untypedDS += len(reg.SupportedDataSources())
	}

	fmt.Printf("typed resource: %d, untyped resource: %d total: %d\n", typedR, untypedR, typedR+untypedR)
	fmt.Printf("typed data: %d, untyped data: %d total: %d\n", typedDS, untypedDS, typedDS+untypedDS)
}

func main() {
	if len(os.Args) <= 1 {
		statistic()
		return
	}
	logic.Run(os.Args[1:]...)
}
