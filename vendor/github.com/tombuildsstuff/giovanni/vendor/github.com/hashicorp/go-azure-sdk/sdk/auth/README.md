# Package: `github.com/hashicorp/go-azure-sdk/sdk/auth`

This package contains Authorizers which can be used to authenticate calls to the Azure APIs for use with `hashicorp/go-azure-sdk`. 

## Example: Authenticating using the Azure CLI

```go
package main

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func main() {
	environment := environments.Public
	credentials := auth.Credentials{
		Environment:                       environment,
		EnableAuthenticatingUsingAzureCLI: true,
	}
	authorizer, err := auth.NewAuthorizerFromCredentials(context.TODO(), credentials, environment.MSGraph)
	if err != nil {
		log.Fatalf("building authorizer from credentials: %+v", err)
	}
	// ...
}
```

## Example: Authenticating using a Client Certificate

```go
package main

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func main() {
	environment := environments.Public
	credentials := auth.Credentials{
		Environment: environment,
		EnableAuthenticatingUsingClientCertificate: true,
		ClientCertificatePath:                      "/path/to/cert.pfx",
		ClientCertificatePassword:                  "somepassword",
	}
	authorizer, err := auth.NewAuthorizerFromCredentials(context.TODO(), credentials, environment.MSGraph)
	if err != nil {
		log.Fatalf("building authorizer from credentials: %+v", err)
	}
	// ..
}
```

## Example: Authenticating using a Client Secret

```go
import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func main() {
	environment := environments.Public
	credentials := auth.Credentials{
		Environment:                           environment,
		EnableAuthenticatingUsingClientSecret: true,
		ClientSecret:                          "some-secret-value",
	}
	authorizer, err := auth.NewAuthorizerFromCredentials(context.TODO(), credentials, environment.MSGraph)
	if err != nil {
		log.Fatalf("building authorizer from credentials: %+v", err)
	}
	// ..
}
```

## Example: Authenticating using a Managed Identity

```go
package main

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func main() {
	environment := environments.Public
	credentials := auth.Credentials{
		Environment:                              environment,
		EnableAuthenticatingUsingManagedIdentity: true,
	}
	authorizer, err := auth.NewAuthorizerFromCredentials(context.TODO(), credentials, environment.MSGraph)
	if err != nil {
		log.Fatalf("building authorizer from credentials: %+v", err)
	}
	// ..
}
```

## Example: Authenticating using GitHub OIDC

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func main() {
	environment := environments.Public
	credentials := auth.Credentials{
		Environment:                         environment,
		EnableAuthenticationUsingGitHubOIDC: true,
		GitHubOIDCTokenRequestURL:           os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"),
		GitHubOIDCTokenRequestToken:         os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN"),
	}
	authorizer, err := auth.NewAuthorizerFromCredentials(context.TODO(), credentials, environment.MSGraph)
	if err != nil {
		log.Fatalf("building authorizer from credentials: %+v", err)
	}
	// ..
}
```

## Example: Authenticating using OIDC

```go
package main

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func main() {
	environment := environments.Public
	credentials := auth.Credentials{
		Environment:                   environment,
		EnableAuthenticationUsingOIDC: true,
		OIDCAssertionToken:            "some-token",
	}
	authorizer, err := auth.NewAuthorizerFromCredentials(context.TODO(), credentials, environment.MSGraph)
	if err != nil {
		log.Fatalf("building authorizer from credentials: %+v", err)
	}
	// ..
}
```