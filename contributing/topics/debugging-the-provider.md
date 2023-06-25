# Debugging the Provider

The provider can be debugged in a number of ways:

- [Adding Log Messages](#logs)
- [Proxying Traffic](#proxy)
- [Attaching a Debugger](#debugger-delve)

## Logs

Adding logging is the most basic, and simplest of ways to debug the provider. Log messages can be added with logging statements such as:

```go

// info message
id, err := parse.SomeResourceId(d.Id())
if err != nil {
	return err
}
log.Printf("[INFO] %s was not found - removing from state", *id)

// debug message
log.Printf("[DEBUG] Importing Resource - parsing %q", d.Id())
```

> **Note:** When logging, lean on the Resource ID Struct (returned from the Parse method above - as shown in the 'info' example above) rather than outputting the Raw Resource ID value (as shown in the debug example above)

These can be viewed by running Terraform (or the Acceptance Test) with logging enabled:

```shell
$ TF_LOG=INFO terraform apply
$ TF_LOG=DEBUG make acctests SERVICE='<service>' TESTARGS='-run=<nameOfTheTest>' TESTTIMEOUT='60m'
```

For more information see [the official Terraform plugin logging documentation](https://www.terraform.io/plugin/log/managing).

## Proxy

A useful step between logging and actual debugging is proxying the traffic through a web debugging proxy such as [Charles Proxy (macOS)](https://www.charlesproxy.com/) or [Fiddler (Windows)](https://www.telerik.com/fiddler). These allow inspection of the web traffic between the provider and Azure to confirm what is actually going across the wire.

You will need to enable HTTPS proxy support (usually by adding a certificate to your system) and then assuming the proxy is running on port `8888`:

```shell
$ http_proxy=http://localhost:8888 https_proxy=http://localhost:8888 terraform apply
$ http_proxy=http://localhost:8888 https_proxy=http://localhost:8888 make acctests SERVICE='<service>' TESTARGS='-run=<nameOfTheTest>' TESTTIMEOUT='60m' 
```

## Debugger (delve)

And finally the most advanced and powerful debugging tool is attaching a debugger such as delve to the provider whilst it is running.

We generally recommend using [Goland](https://jetbrains.com/go) as it provides (amongst other features) native integrations for debugging - see [OpenCredo's blog post](https://opencredo.com/blogs/running-a-terraform-provider-with-a-debugger/) for an example - however it's also possible to use [VSCode and the delve CLI](https://www.terraform.io/plugin/debugging) - configuring these is outside of the scope of this project.
