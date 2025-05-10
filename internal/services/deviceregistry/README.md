Overview of Device Registry Acceptance Tests
===
The Azure Device Registry service has several arc-enabled resources, including Assets and Asset Endpoint Profiles (AEPs). These resources will only create successfully if there is an arc-enabled Kubernetes Cluster in Azure that runs all of Azure IoT Operations' (AIO) service and has a corresponding Custom Location. You can learn more about AIO [here](https://learn.microsoft.com/en-us/azure/iot-operations/). Because of this requirement, this makes acceptance testing Assets and AEP resources more complex since the tests must setup an AIO cluster.

The solution that we have adapts that from the [Custom Location tests](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/internal/services/extendedlocation/extended_location_custom_location_test.go), but it is a little more complex. Here is an overview of the process the Device Registry tests need to do:
1. First, each of the Device Registry acceptance tests apply a Terraform template to create an Azure Linux VM and all of the VM's infrastructure resources (e.g. public IP address, subnet, the resource group that will hold everything, etc). The VM will host the AIO cluster. The tests will also provision a bash script file to the VM which will execute all the commands needed to setup the AIO cluster. The bash script Terraform template file is [setup_aio_cluster.sh.tftpl](./testdata/setup_aio_cluster.sh.tftpl) and can be found in the `testdata` directory. The tests do not run the bash script yet.

2. Before the Assets/AEPs resources are created, a `PreConfig` step is run. The tests execute some Go code to fetch the VM's public IP address and then uses the IP address to SSH into the VM and execute the bash script on the VM. The bash script will install Azure CLI and setup a [K3s cluster](https://k3s.io/) on the VM, and then run the AZ CLI commands from this [AIO quickstart](https://learn.microsoft.com/en-us/azure/iot-operations/get-started-end-to-end-sample/quickstart-deploy) to arc-enable the cluster and setup AIO services on it (which will also create the Custom Location).
    - We must do it this way because even with the `depends_on` property, the tests do not wait for the VM to finish its `remote-exec` to run the bash script. Thus, the tests will fail as they will try to create Assets/AEPs while the AIO cluster and Custom Location are provisioning, throwing a "Custom Location (or other AIO resource) does not exist" error. This is the only way to sequentially execute the bash script to setup the AIO cluster and block the tests from prematurely creating the Asset/AEPs, as attempts to use `null_resource`, Go's `time.sleep()`, etc ended up not working (and stopped `remote-exec` from completing). Also, setting a time limit to wait for the cluster to finish is not recommended as the time to finish script execution can take anywhere between 2000-3500 seconds or even more.

3. Once the bash script completes execution, then the rest of the test proceeds as normal; the test creates the Asset/AEP on the AIO cluster. Note: each resource's test scenarios currently creates a separate AIO cluster for each test scenario. So please make sure the Azure subscription has enough resources to concurrently create multiple VMs. 

4. When a test scenario finishes, the cleanup steps will run. The VM, Asset/AEP resource, and other VM infra resources will automatically be destroyed by the test cleanup. However, the AIO cluster and its own resources were created by the VM, not Terraform, so they would not get targeted for deletion by the tests. Fortunately, the tests created the resource group that contains all of these resources. So we specify to the acceptance tests to delete the entire resource group to cleanup the AIO cluster resources, as well. That is why the `prevent_deletion_if_contains_resources` flag is set to false in the tests:
```
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}
```

How to run Device Registry Acceptance Tests
===
1. On your own machine, login to Azure CLI as a user with ownership permissions of the Azure subscription the acceptance tests will run on.

2. Run the following commands to enable the providers in your Azure subscription so that the AIO Cluster setup steps will not fail. You only have to do this once for your subscription and after that you can skip this step.
```bash
az provider register -n "Microsoft.ExtendedLocation"
az provider register -n "Microsoft.Kubernetes"
az provider register -n "Microsoft.KubernetesConfiguration"
az provider register -n "Microsoft.IoTOperations"
az provider register -n "Microsoft.DeviceRegistry"
az provider register -n "Microsoft.SecretSyncController"
```

3. Run `az ad sp show --id bc313c14-388c-4e7d-a58e-70017303ee3b --query id -o tsv` to get the Custom Location RP's Entra App Object ID. Store it in an environment variable `ARM_ENTRA_APP_OBJECT_ID` (`export ARM_ENTRA_APP_OBJECT_ID=<object ID>`). In theory, you only need to run the `az ad sp show` command once because once you have the object ID, you can reuse that object ID in the acceptance test pipeline for future test runs.

4. The following environment variables need to be set to run the Acceptance Tests. Make sure that the Service Principal running the tests has ownership permissions of the subscription so that the Azure CLI commands in the setup script do not fail.
```bash
# ID of the Azure subscription that the acceptance tests will run on
export ARM_SUBSCRIPTION_ID=<subscription ID> 

# The Client ID of the Service Principal that will run the acceptance tests.
export ARM_CLIENT_ID=<client ID>

# The password of the Service Principal that will run the acceptance tests.
export ARM_CLIENT_SECRET=<client secret> 

# The Object ID of the Custom Locations RP's Entra App, as mentioned in previous step.
export ARM_ENTRA_APP_OBJECT_ID=<object ID> 
```

5. Run the acceptance tests as normal.
