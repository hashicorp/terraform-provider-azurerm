package main

import (
	"flag"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func main() {
	isResource := flag.Bool("resource", true, "True if this is a Resource, False if this is a Data Source")
	resourceName := flag.String("name", "", "The name of the Data Source/Resource which should be generated")
	flag.Parse()

	if err := run(*isResource, *resourceName); err != nil {
		panic(err)
	}
}

func run(isResource bool, resourceName string) error {
	if isResource {
		return runForResource(resourceName)
	}

	return runForDataSource(resourceName)
}

func runForDataSource(resourceName string) error {
	var dataSource *schema.Resource

	for _, service := range provider.SupportedServices() {
		for key, ds := range service.SupportedDataSources() {
			if key == resourceName {
				dataSource = ds
				break
			}
		}
	}

	if dataSource == nil {
		return fmt.Errorf("Data Source %q was not registered!", resourceName)
	}

	return buildDocumentationForDataSource(resourceName, dataSource)
}

func runForResource(resourceName string) error {
	var resource *schema.Resource

	for _, service := range provider.SupportedServices() {
		for key, rs := range service.SupportedResources() {
			if key == resourceName {
				resource = rs
				break
			}
		}
	}

	if resource == nil {
		return fmt.Errorf("Resource %q was not registered!", resourceName)
	}

	return buildDocumentationForResource(resourceName, resource)
}

func buildDocumentationForDataSource(resourceName string, resource *schema.Resource) error {
	return nil
}

func buildDocumentationForResource(resourceName string, resource *schema.Resource) error {
	return nil
}
