// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure_ssis

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeAzureSsisRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName)

	resp, err := client.Get(ctx, *id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IntegrationRuntimeName)
	d.Set("data_factory_id", dataFactoryId.ID())

	if model := resp.Model; model != nil {
		runTime, ok := model.Properties.(integrationruntimes.ManagedIntegrationRuntime)
		if !ok {
			return fmt.Errorf("asserting `IntegrationRuntime` as `ManagedIntegrationRuntime` for %s", *id)
		}

		d.Set("description", runTime.Description)

		if computeProps := runTime.TypeProperties.ComputeProperties; computeProps != nil {
			d.Set("location", location.NormalizeNilable(computeProps.Location))
			d.Set("node_size", computeProps.NodeSize)
			d.Set("number_of_nodes", computeProps.NumberOfNodes)
			d.Set("max_parallel_executions_per_node", computeProps.MaxParallelExecutionsPerNode)

			if err := d.Set("vnet_integration", flattenDataFactoryIntegrationRuntimeAzureSsisVnetIntegration(computeProps.VNetProperties)); err != nil {
				return fmt.Errorf("setting `vnet_integration`: %+v", err)
			}

			if err := d.Set("copy_compute_scale", flattenDataFactoryIntegrationRuntimeAzureSsisCopyComputeScale(computeProps.CopyComputeScaleProperties)); err != nil {
				return fmt.Errorf("setting `copy_compute_scale`: %+v", err)
			}

			if err := d.Set("pipeline_external_compute_scale", flattenDataFactoryIntegrationRuntimeAzureSsisPipelineExternalComputeScaleProperties(computeProps.PipelineExternalComputeScaleProperties)); err != nil {
				return fmt.Errorf("setting `pipeline_external_compute_scale`: %+v", err)
			}
		}

		if ssisProps := runTime.TypeProperties.SsisProperties; ssisProps != nil {
			d.Set("edition", pointer.FromEnum(ssisProps.Edition))
			d.Set("license_type", pointer.FromEnum(ssisProps.LicenseType))

			if err := d.Set("catalog_info", flattenDataFactoryIntegrationRuntimeAzureSsisCatalogInfo(ssisProps.CatalogInfo, d)); err != nil {
				return fmt.Errorf("setting `catalog_info`: %+v", err)
			}

			if err := d.Set("credential_name", flattenDataFactoryIntegrationRuntimeUserAssignedCredential(ssisProps.Credential)); err != nil {
				return fmt.Errorf("setting `credential_name`: %+v", err)
			}

			if err := d.Set("custom_setup_script", flattenDataFactoryIntegrationRuntimeAzureSsisCustomSetupScript(ssisProps.CustomSetupScriptProperties, d)); err != nil {
				return fmt.Errorf("setting `custom_setup_script`: %+v", err)
			}

			if err := d.Set("express_custom_setup", flattenDataFactoryIntegrationRuntimeAzureSsisExpressCustomSetUp(ssisProps.ExpressCustomSetupProperties, d)); err != nil {
				return fmt.Errorf("setting `express_custom_setup`: %+v", err)
			}

			if err := d.Set("package_store", flattenDataFactoryIntegrationRuntimeAzureSsisPackageStore(ssisProps.PackageStores)); err != nil {
				return fmt.Errorf("setting `package_store`: %+v", err)
			}

			if err := d.Set("proxy", flattenDataFactoryIntegrationRuntimeAzureSsisProxy(ssisProps.DataProxyProperties)); err != nil {
				return fmt.Errorf("setting `proxy`: %+v", err)
			}
		}

		if err := d.Set("express_vnet_integration", flattenDataFactoryIntegrationRuntimeCustomerVnetIntegration(runTime.TypeProperties.CustomerVirtualNetwork)); err != nil {
			return fmt.Errorf("setting `express_vnet_integration`: %+v", err)
		}
	}

	return nil
}

func readBackSensitiveValue(input []any, propertyName string, filters map[string]string) string {
	if len(input) == 0 {
		return ""
	}
	for _, item := range input {
		raw := item.(map[string]any)
		found := true
		for k, v := range filters {
			if raw[k].(string) != v {
				found = false
				break
			}
		}
		if found {
			return raw[propertyName].(string)
		}
	}
	return ""
}

func formatDataFactoryIntegrationRuntimeElasticPool(input string) string {
	return fmt.Sprintf(`ELASTIC_POOL(name="%s")`, input)
}

func parseDataFactoryIntegrationRuntimeElasticPool(input string) (string, bool) {
	matches := regexp.MustCompile(`^ELASTIC_POOL\(name="(.+)"\)$`).FindStringSubmatch(input)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}
