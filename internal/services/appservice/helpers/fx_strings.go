package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func decodeApplicationStackLinux(fxString string) ApplicationStackLinux {
	parts := strings.Split(fxString, "|")
	result := ApplicationStackLinux{}
	if len(parts) != 2 {
		return result
	}

	switch parts[0] {
	case "DOTNETCORE", "DOTNET":
		result.NetFrameworkVersion = parts[1]

	case "NODE":
		result.NodeVersion = parts[1]

	case "JAVA", "TOMCAT", "JBOSSEAP":
		result.JavaServer = parts[0]
		javaParts := strings.Split(parts[1], "-")
		if len(javaParts) == 2 {
			// e.g. 8-jre8
			result.JavaServerVersion = javaParts[0]
			result.JavaVersion = javaParts[1]
		} else {
			// e.g. 8u242 or 11.0.9
			result.JavaVersion = parts[1]
		}

	case "PHP":
		result.PhpVersion = parts[1]

	case "PYTHON":
		result.PythonVersion = parts[1]

	case "RUBY":
		result.RubyVersion = parts[1]

	default: // DOCKER is the expected default here as "custom" images require it
		if dockerParts := strings.Split(parts[1], ":"); len(dockerParts) == 2 {
			result.DockerImage = dockerParts[0]
			result.DockerImageTag = dockerParts[1]
		}
	}

	return result
}

func EncodeFunctionAppLinuxFxVersion(input []ApplicationStackLinuxFunctionApp) *string {
	if len(input) == 0 {
		return utils.String("")
	}

	appStack := input[0]
	var appType, appString string
	switch {
	case appStack.NodeVersion != "":
		appType = "Node"
		appString = appStack.NodeVersion
	case appStack.DotNetVersion != "":
		appType = "DotNet"
		appString = appStack.DotNetVersion
	case appStack.PythonVersion != "":
		appType = "Python"
		appString = appStack.PythonVersion
	case appStack.JavaVersion != "":
		appType = "Java"
		appString = appStack.JavaVersion
	case len(appStack.Docker) > 0 && appStack.Docker[0].ImageName != "":
		appType = "Docker"
		dockerCfg := appStack.Docker[0]
		if dockerCfg.RegistryURL != "" {
			appString = fmt.Sprintf("%s/%s:%s", strings.Trim(dockerCfg.RegistryURL, "/"), dockerCfg.ImageName, dockerCfg.ImageTag)
		} else {
			appString = fmt.Sprintf("%s:%s", dockerCfg.ImageName, dockerCfg.ImageTag)
		}
	}

	return utils.String(fmt.Sprintf("%s|%s", appType, appString))
}

func DecodeFunctionAppLinuxFxVersion(input string) ([]ApplicationStackLinuxFunctionApp, error) {
	if input == "" {
		// This is a valid string for "Custom" stack which we picked up earlier, so we can skip here
		return nil, nil
	}

	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unrecognised LinuxFxVersion format received, got %s", input)
	}

	result := make([]ApplicationStackLinuxFunctionApp, 0)

	switch strings.ToLower(parts[0]) {
	case "dotnet":
		appStack := ApplicationStackLinuxFunctionApp{DotNetVersion: parts[1]}
		result = append(result, appStack)

	case "node":
		appStack := ApplicationStackLinuxFunctionApp{NodeVersion: parts[1]}
		result = append(result, appStack)

	case "python":
		appStack := ApplicationStackLinuxFunctionApp{PythonVersion: parts[1]}
		result = append(result, appStack)

	case "java":
		appStack := ApplicationStackLinuxFunctionApp{JavaVersion: parts[1]}
		result = append(result, appStack)

	case "docker":
		docker, err := decodeFunctionAppDockerFxString(parts[1])
		if err != nil {
			return nil, err
		}
		appStack := ApplicationStackLinuxFunctionApp{Docker: docker}
		result = append(result, appStack)
	}

	return result, nil
}

func decodeFunctionAppDockerFxString(input string) ([]ApplicationStackDocker, error) {

	return nil, nil // TODO - implement this
}
