## Tool: `update-api-version`

This tool automatically updates import references from one API version to another API version when using `hashicorp/go-azure-sdk`.

This works against a Service and takes both the old and the new API versions, for example:

```sh
./update-api-version -service="advisor" -old-api-version="2020-01-01" -new-api-version="2022-01-01"
```

The arguments are:

* `service` - the name of the Service Package (e.g. `compute`, `network` etc).
* `old-api-version` - the existing API version (for example `2020-01-01` or `2020-01-01-preview`) which should be replaced.
* `new-api-version` - the new API version which should be used in place of the value for `old-api-version`.

---

This tool supports updating both explicit SDK references, that is:

```
import (
    "github.com/hashicorp/go-azure-sdk/resource-manager/{service}/{api-version}/{sdk}"
)
```

as well as Meta Clients (which includes updating the import alias):

```
import (
    service_api_version "github.com/hashicorp/go-azure-sdk/resource-manager/{service}/{api-version}"
)
```
