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

	switch strings.ToUpper(parts[0]) {
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
	if len(input) == 0 || input[0].CustomHandler {
		return utils.String("")
	}

	appStack := input[0]
	var appType, appString string
	switch {
	case appStack.NodeVersion != "":
		appType = "NODE"
		appString = appStack.NodeVersion

	case appStack.DotNetVersion != "":
		if appStack.DotNetIsolated {
			appType = "DOTNET-ISOLATED"
		} else {
			appType = "DOTNET"
		}
		appString = appStack.DotNetVersion

	case appStack.PythonVersion != "":
		appType = "PYTHON"
		appString = appStack.PythonVersion

	case appStack.JavaVersion != "":
		appType = "JAVA"
		appString = appStack.JavaVersion

	case appStack.PowerShellCoreVersion != "":
		appType = "POWERSHELL"
		appString = appStack.PowerShellCoreVersion

	case len(appStack.Docker) > 0 && appStack.Docker[0].ImageName != "":
		appType = "DOCKER"
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

	case "dotnet-isolated":
		appStack := ApplicationStackLinuxFunctionApp{DotNetVersion: parts[1], DotNetIsolated: true}
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

	case "powershell":
		appStack := ApplicationStackLinuxFunctionApp{PowerShellCoreVersion: parts[1]}
		result = append(result, appStack)

	case "docker":
		// This is handled as part of unpacking the app_settings using DecodeFunctionAppDockerFxString but included here for signposting as this is not intuitive.
	}

	return result, nil
}

func DecodeFunctionAppDockerFxString(input string, partial ApplicationStackDocker) ([]ApplicationStackDocker, error) {
	if input == "" {
		// This is a valid string for "Custom" stack which we picked up earlier, so we can skip here
		return nil, nil
	}

	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unrecognised LinuxFxVersion format received, got %s", input)
	}

	if !strings.EqualFold(parts[0], "docker") {
		return nil, fmt.Errorf("expected a docker FX version, got %q", parts[0])
	}

	dockerParts := strings.Split(strings.TrimPrefix(parts[1], partial.RegistryURL), ":")
	if len(dockerParts) != 2 {
		return nil, fmt.Errorf("invalid docker image reference %q", parts[1])
	}

	partial.ImageName = strings.TrimPrefix(dockerParts[0], "/")
	partial.ImageTag = dockerParts[1]

	return []ApplicationStackDocker{partial}, nil
}

func EncodeFunctionAppWindowsFxVersion(input []ApplicationStackWindowsFunctionApp) *string {
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
		if appStack.DotNetIsolated {
			appType = "DotNet-Isolated"
		} else {
			appType = "DotNet"
		}
		appString = appStack.DotNetVersion

	case appStack.JavaVersion != "":
		appType = "Java"
		appString = appStack.JavaVersion

	case appStack.PowerShellCoreVersion != "":
		appType = "PowerShell"
		appString = appStack.PowerShellCoreVersion
	}

	return utils.String(fmt.Sprintf("%s|%s", appType, appString))
}

func DecodeFunctionAppWindowsFxVersion(input string) ([]ApplicationStackWindowsFunctionApp, error) {
	if input == "" {
		// This is a valid string for "Custom" stack which we picked up earlier, so we can skip here
		return nil, nil
	}

	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unrecognised WindowsFxVersion format received, got %s", input)
	}

	result := make([]ApplicationStackWindowsFunctionApp, 0)

	switch strings.ToLower(parts[0]) {
	case "dotnet":
		appStack := ApplicationStackWindowsFunctionApp{DotNetVersion: parts[1]}
		result = append(result, appStack)

	case "dotnet-isolated":
		appStack := ApplicationStackWindowsFunctionApp{DotNetVersion: parts[1], DotNetIsolated: true}
		result = append(result, appStack)

	case "node":
		appStack := ApplicationStackWindowsFunctionApp{NodeVersion: parts[1]}
		result = append(result, appStack)

	case "java":
		appStack := ApplicationStackWindowsFunctionApp{JavaVersion: parts[1]}
		result = append(result, appStack)

	case "powershell":
		appStack := ApplicationStackWindowsFunctionApp{PowerShellCoreVersion: parts[1]}
		result = append(result, appStack)
	}

	return result, nil
}

func JavaLinuxFxStringBuilder(javaMajorVersion, javaServer, javaServerVersion string) string {
	switch javaMajorVersion {
	case "8":
		{
			switch javaServer {
			case LinuxJavaServerJava:
				if strings.Contains(javaServerVersion, "u") {
					return fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion) // e.g. JAVA|8u302
				} else {
					return fmt.Sprintf("%s|%s-jre8", LinuxJavaServerJava, javaServerVersion) // e.g. "JAVA|8-jre8"
				}
			case LinuxJavaServerTomcat:
				if len(strings.Split(javaServerVersion, ".")) == 3 {
					return fmt.Sprintf("%s|%s-java8", LinuxJavaServerTomcat, javaServerVersion) // e.g. TOMCAT|10.0.20-java8
				} else {
					return fmt.Sprintf("%s|%s-jre8", LinuxJavaServerTomcat, javaServerVersion) // e.g. TOMCAT|10.0-jre8
				}
			case LinuxJavaServerJboss:
				return fmt.Sprintf("%s|%s-java8", LinuxJavaServerJboss, javaServerVersion)
			}
		}
	case "11":
		switch javaServer {
		case LinuxJavaServerJava:
			if len(strings.Split(javaServerVersion, ".")) == 3 {
				return fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion) // e.g. JAVA|11.0.13
			} else {
				return fmt.Sprintf("%s|%s-java11", LinuxJavaServerJava, javaServerVersion) // e.g.JAVA|11-java1
			}
		case LinuxJavaServerTomcat:
			return fmt.Sprintf("%s|%s-java11", LinuxJavaServerTomcat, javaServerVersion) // e.g. TOMCAT|10.0-java11 and TOMCAT|10.0.20-java11

		case LinuxJavaServerJboss:
			return fmt.Sprintf("%s|%s-java11", LinuxJavaServerJboss, javaServerVersion) // e.g. TOMCAT|10.0-java11 and TOMCAT|10.0.20-java11// e.g. JBOSSEAP|7-java11 / JBOSSEAP|7.4.2-java11
		}

	case "17":
		switch javaServer {
		case LinuxJavaServerJava:
			if len(strings.Split(javaServerVersion, ".")) == 3 {
				return fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion) // "JAVA|17.0.2"
			} else {
				return fmt.Sprintf("%s|%s-java17", LinuxJavaServerJava, javaServerVersion) // "JAVA|17-java17"
			}

		case LinuxJavaServerTomcat:
			return fmt.Sprintf("%s|%s-java17", LinuxJavaServerTomcat, javaServerVersion) // e,g, TOMCAT|10.0-java17 / TOMCAT|10.0.20-java17
		default:
			return fmt.Sprintf("%s|%s-java17", javaServer, javaServerVersion)
		}

	default:
		return fmt.Sprintf("%s|%s-%s", javaServer, javaServerVersion, javaMajorVersion)

	}
	return ""
}
