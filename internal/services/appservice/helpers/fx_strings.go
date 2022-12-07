package helpers

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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

	case "GO":
		result.GoVersion = parts[1]

	case "NODE":
		result.NodeVersion = parts[1]

	case LinuxJavaServerJava:
		result.JavaServer = LinuxJavaServerJava
		javaParts := strings.Split(parts[1], "-")
		if strings.HasPrefix(parts[1], "8") {
			result.JavaVersion = "8"
		}
		if strings.HasPrefix(javaParts[0], "11") {
			result.JavaVersion = "11"
		}
		if strings.HasPrefix(javaParts[0], "17") {
			result.JavaVersion = "17"
		}
		result.JavaServerVersion = javaParts[0]

	case LinuxJavaServerTomcat:
		result.JavaServer = LinuxJavaServerTomcat
		javaParts := strings.Split(parts[1], "-")
		if len(javaParts) == 2 {
			result.JavaServerVersion = javaParts[0]
			javaVersion := strings.TrimPrefix(javaParts[1], "jre")
			javaVersion = strings.TrimPrefix(javaVersion, "java")
			result.JavaVersion = javaVersion
		}

	case LinuxJavaServerJboss:
		result.JavaServer = LinuxJavaServerJboss
		javaParts := strings.Split(parts[1], "-")
		if len(javaParts) == 2 {
			result.JavaServerVersion = javaParts[0]
			javaVersion := strings.TrimPrefix(javaParts[1], "jre")
			javaVersion = strings.TrimPrefix(javaVersion, "java")
			result.JavaVersion = javaVersion
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

func JavaLinuxFxStringBuilder(javaMajorVersion, javaServer, javaServerVersion string) (*string, error) {
	switch javaMajorVersion {
	case "8":
		{
			switch javaServer {
			case LinuxJavaServerJava:
				if strings.Contains(javaServerVersion, "u") {
					return pointer.To(fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion)), nil // e.g. JAVA|8u302
				} else {
					return pointer.To(fmt.Sprintf("%s|%s-jre8", LinuxJavaServerJava, javaServerVersion)), nil // e.g. "JAVA|8-jre8"
				}
			case LinuxJavaServerTomcat:
				if len(strings.Split(javaServerVersion, ".")) == 3 {
					return pointer.To(fmt.Sprintf("%s|%s-java8", LinuxJavaServerTomcat, javaServerVersion)), nil // e.g. TOMCAT|10.0.20-java8
				} else {
					return pointer.To(fmt.Sprintf("%s|%s-jre8", LinuxJavaServerTomcat, javaServerVersion)), nil // e.g. TOMCAT|10.0-jre8
				}
			case LinuxJavaServerJboss:
				return pointer.To(fmt.Sprintf("%s|%s-java8", LinuxJavaServerJboss, javaServerVersion)), nil
			}
		}
	case "11":
		switch javaServer {
		case LinuxJavaServerJava:
			if len(strings.Split(javaServerVersion, ".")) == 3 {
				return pointer.To(fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion)), nil // e.g. JAVA|11.0.13
			} else {
				return pointer.To(fmt.Sprintf("%s|%s-java11", LinuxJavaServerJava, javaServerVersion)), nil // e.g.JAVA|11-java1
			}
		case LinuxJavaServerTomcat:
			return pointer.To(fmt.Sprintf("%s|%s-java11", LinuxJavaServerTomcat, javaServerVersion)), nil // e.g. TOMCAT|10.0-java11 and TOMCAT|10.0.20-java11

		case LinuxJavaServerJboss:
			return pointer.To(fmt.Sprintf("%s|%s-java11", LinuxJavaServerJboss, javaServerVersion)), nil // e.g. TOMCAT|10.0-java11 and TOMCAT|10.0.20-java11// e.g. JBOSSEAP|7-java11 / JBOSSEAP|7.4.2-java11
		}

	case "17":
		switch javaServer {
		case LinuxJavaServerJava:
			if len(strings.Split(javaServerVersion, ".")) == 3 {
				return pointer.To(fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion)), nil // "JAVA|17.0.2"
			} else {
				return pointer.To(fmt.Sprintf("%s|%s-java17", LinuxJavaServerJava, javaServerVersion)), nil // "JAVA|17-java17"
			}

		case LinuxJavaServerTomcat:
			return pointer.To(fmt.Sprintf("%s|%s-java17", LinuxJavaServerTomcat, javaServerVersion)), nil // e,g, TOMCAT|10.0-java17 / TOMCAT|10.0.20-java17
		case LinuxJavaServerJboss:
			return nil, fmt.Errorf("java 17 is not supported on %s", LinuxJavaServerJboss)
		default:
			return pointer.To(fmt.Sprintf("%s|%s-java17", javaServer, javaServerVersion)), nil
		}

	default:
		return pointer.To(fmt.Sprintf("%s|%s-%s", javaServer, javaServerVersion, javaMajorVersion)), nil

	}
	return nil, fmt.Errorf("unsupported combination of `java_version`, `java_server`, and `java_server_version`")
}
