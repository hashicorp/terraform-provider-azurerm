package main

import (
	// "fmt"
	"log"
	"context"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	// "encoding/json"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
)

func main() {
	// create()	
	delete()
}

func create() {
	// client := network.NewBastionHostsClient("4b89b857-a820-4f5e-b02a-da6d5e180752")
	

	log.Println("[INFO] preparing arguments for Azure Bastion Host creation.")

	// resourceGroup := "test"
	// name := "test"
	location := "westeurope"
	dnsName := "testasfqwf"
					subID := "subid"
	pipID := "pip"

	subnetID := network.SubResource{
		ID: &subID,
	}

	publicIPAddressID := network.SubResource{
		ID: &pipID,
	}

	bastionHostIPConfigurationPropertiesFormat := network.BastionHostIPConfigurationPropertiesFormat{
		Subnet:          &subnetID,
		PublicIPAddress: &publicIPAddressID,
	}

	bastionHostIPConfiguration := []network.BastionHostIPConfiguration{{
		Name: to.StringPtr("Name"),
		BastionHostIPConfigurationPropertiesFormat: &bastionHostIPConfigurationPropertiesFormat,
	},
	}

	properties := network.BastionHostPropertiesFormat{
		IPConfigurations: &bastionHostIPConfiguration,
		DNSName:          &dnsName,
	}

	parameters := network.BastionHost{
		Location:                    &location,
		BastionHostPropertiesFormat: &properties,
	}

	j, _ := parameters.MarshalJSON()
	log.Println(string(j))
	// // // creation
	// // future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	// // if err != nil {
	// 	// return fmt.Errorf("Error creating/updating Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	// // 	return err
	// // }
	// // if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
	// // 	return fmt.Errorf("Error waiting for creation/update of Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	// // }
	// // // check presence
	// // read, err := client.Get(ctx, resourceGroup, name)
	// // if err != nil {
	// // 	return fmt.Errorf("Error retrieving Bastion Host %q (Resource Group %q): %+v", name, resourceGroup, err)
	// // }
}

func delete() {
	// create a VirtualNetworks client
	client := network.NewBastionHostsClient("4b89b857-a820-4f5e-b02a-da6d5e180752")
	ctx := context.Background()

	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err == nil {
		client.Authorizer = authorizer
	}
	
	s, err := client.Delete(ctx, "example-resources-2", "testbastion")
	if err != nil {
		log.Println(err)
	}
	log.Println(s)
}