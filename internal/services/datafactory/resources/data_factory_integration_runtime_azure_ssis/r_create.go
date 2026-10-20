// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure_ssis

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeAzureSsisCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := integrationruntimes.NewIntegrationRuntimeID(dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, integrationruntimes.DefaultGetOperationOptions())
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_azure_ssis", id.ID())
		}
	}

	managedIntegrationRuntime := integrationruntimes.ManagedIntegrationRuntime{
		Description: pointer.To(d.Get("description").(string)),
		Type:        integrationruntimes.IntegrationRuntimeTypeManaged,
		TypeProperties: integrationruntimes.ManagedIntegrationRuntimeTypeProperties{
			ComputeProperties: &integrationruntimes.IntegrationRuntimeComputeProperties{
				Location:                               pointer.To(location.Normalize(d.Get("location").(string))),
				NodeSize:                               pointer.To(d.Get("node_size").(string)),
				NumberOfNodes:                          pointer.To(int64(d.Get("number_of_nodes").(int))),
				MaxParallelExecutionsPerNode:           pointer.To(int64(d.Get("max_parallel_executions_per_node").(int))),
				VNetProperties:                         expandDataFactoryIntegrationRuntimeAzureSsisVirtualNetwork(d.Get("vnet_integration").([]any)),
				CopyComputeScaleProperties:             expandDataFactoryIntegrationRuntimeAzureSsisCopyComputeScale(d.Get("copy_compute_scale").([]any)),
				PipelineExternalComputeScaleProperties: expandDataFactoryIntegrationRuntimeAzureSsisPipelineExternalComputeScale(d.Get("pipeline_external_compute_scale").([]any)),
			},
			SsisProperties: &integrationruntimes.IntegrationRuntimeSsisProperties{
				LicenseType:                  pointer.ToEnum[integrationruntimes.IntegrationRuntimeLicenseType](d.Get("license_type").(string)),
				DataProxyProperties:          expandDataFactoryIntegrationRuntimeAzureSsisProxy(d.Get("proxy").([]any)),
				Edition:                      pointer.ToEnum[integrationruntimes.IntegrationRuntimeEdition](d.Get("edition").(string)),
				ExpressCustomSetupProperties: expandDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(d.Get("express_custom_setup").([]any)),
				PackageStores:                expandDataFactoryIntegrationRuntimeAzureSsisPackageStore(d.Get("package_store").([]any)),
				Credential:                   expandDataFactoryIntegrationRuntimeAzureSsisCredential(d.Get("credential_name").(string)),
				CatalogInfo:                  expandDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(d.Get("catalog_info").([]any)),
				CustomSetupScriptProperties:  expandDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(d.Get("custom_setup_script").([]any)),
			},
			CustomerVirtualNetwork: expandDataFactoryIntegrationRuntimeCustomerVirtualNetwork(d.Get("express_vnet_integration").([]any)),
		},
	}

	integrationRuntime := integrationruntimes.IntegrationRuntimeResource{
		Name:       pointer.To(id.IntegrationRuntimeName),
		Properties: managedIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, id, integrationRuntime, integrationruntimes.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryIntegrationRuntimeAzureSsisRead(d, meta)
}
