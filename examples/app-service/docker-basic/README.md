# Example: a Linux App Service running a Docker container

This example provisions a Linux App Service which runs a single Docker container.

## Variables

- `prefix` - (Required) The prefix used for all resources in this example.
- `location` - (Required) Azure Region in which all resources in this example should be provisioned.

## Outputs

- `app_name` - The name of the app.
- `app_url` - The default URL to access the app.

## Notes

- The Container is launched on the first HTTP Request, which can take a while.
- Continuous Deployment of a single Docker Container can be achieved using the App Setting `DOCKER_ENABLE_CI` to `true`.
