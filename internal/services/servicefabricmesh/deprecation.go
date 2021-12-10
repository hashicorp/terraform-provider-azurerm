package servicefabricmesh

import "fmt"

func deprecationMessage(resourceName string) string {
	return fmt.Sprintf(`
Service Fabric Mesh is been retired on 2021-04-28 and new clusters can no longer be provisioned.

As such the %q resource is deprecated and will be removed
in version 3.0 of the Azure Provider. 

More information on the retirement can be found here:
https://azure.microsoft.com/en-us/updates/azure-service-fabric-mesh-preview-retirement/`, resourceName)
}
