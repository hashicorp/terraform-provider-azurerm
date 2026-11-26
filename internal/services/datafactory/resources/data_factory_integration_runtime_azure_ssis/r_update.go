// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure_ssis

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeAzureSsisUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `Model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `Properties` was nil", id)
	}

	props, ok := existing.Model.Properties.(integrationruntimes.ManagedIntegrationRuntime)
	if !ok {
		return fmt.Errorf("retrieving %s: asserting `IntegrationRuntime` as `ManagedIntegrationRuntime`", id)
	}

	if props.TypeProperties.ComputeProperties == nil {
		props.TypeProperties.ComputeProperties = &integrationruntimes.IntegrationRuntimeComputeProperties{}
	}
	computeProps := props.TypeProperties.ComputeProperties

	if props.TypeProperties.SsisProperties == nil {
		props.TypeProperties.SsisProperties = &integrationruntimes.IntegrationRuntimeSsisProperties{}
	}
	ssisProps := props.TypeProperties.SsisProperties

	if d.HasChange("node_size") {
		computeProps.NodeSize = pointer.To(d.Get("node_size").(string))
	}

	if d.HasChange("description") {
		props.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("number_of_nodes") {
		computeProps.NumberOfNodes = pointer.To(int64(d.Get("number_of_nodes").(int)))
	}

	if d.HasChange("max_parallel_executions_per_node") {
		computeProps.MaxParallelExecutionsPerNode = pointer.To(int64(d.Get("max_parallel_executions_per_node").(int)))
	}

	if d.HasChange("credential_name") {
		ssisProps.Credential = expandDataFactoryIntegrationRuntimeAzureSsisCredential(d.Get("credential_name").(string))
	}

	if d.HasChange("edition") {
		ssisProps.Edition = pointer.ToEnum[integrationruntimes.IntegrationRuntimeEdition](d.Get("edition").(string))
	}

	if d.HasChange("copy_compute_scale") {
		computeProps.CopyComputeScaleProperties = expandDataFactoryIntegrationRuntimeAzureSsisCopyComputeScale(d.Get("copy_compute_scale").([]any))
	}

	if d.HasChange("express_vnet_integration") {
		props.TypeProperties.CustomerVirtualNetwork = expandDataFactoryIntegrationRuntimeCustomerVirtualNetwork(d.Get("express_vnet_integration").([]any))
	}

	if d.HasChange("license_type") {
		ssisProps.LicenseType = pointer.ToEnum[integrationruntimes.IntegrationRuntimeLicenseType](d.Get("license_type").(string))
	}

	if d.HasChange("vnet_integration") {
		computeProps.VNetProperties = expandDataFactoryIntegrationRuntimeAzureSsisVirtualNetwork(d.Get("vnet_integration").([]any))
	}

	if d.HasChange("custom_setup_script") {
		ssisProps.CustomSetupScriptProperties = expandDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(d.Get("custom_setup_script").([]any))
	}

	if d.HasChange("catalog_info") {
		ssisProps.CatalogInfo = expandDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(d.Get("catalog_info").([]any))
	}

	if d.HasChange("express_custom_setup") {
		ssisProps.ExpressCustomSetupProperties = expandDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(d.Get("express_custom_setup").([]any))
	}

	if d.HasChange("package_store") {
		ssisProps.PackageStores = expandDataFactoryIntegrationRuntimeAzureSsisPackageStore(d.Get("package_store").([]any))
	}

	if d.HasChange("pipeline_external_compute_scale") {
		computeProps.PipelineExternalComputeScaleProperties = expandDataFactoryIntegrationRuntimeAzureSsisPipelineExternalComputeScale(d.Get("pipeline_external_compute_scale").([]any))
	}

	if d.HasChange("proxy") {
		ssisProps.DataProxyProperties = expandDataFactoryIntegrationRuntimeAzureSsisProxy(d.Get("proxy").([]any))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model, integrationruntimes.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDataFactoryIntegrationRuntimeAzureSsisRead(d, meta)
}
