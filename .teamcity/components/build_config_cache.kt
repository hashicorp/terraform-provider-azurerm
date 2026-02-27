import jetbrains.buildServer.configs.kotlin.AbsoluteId
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.buildFeatures.BuildCacheFeature
import jetbrains.buildServer.configs.kotlin.buildSteps.ScriptBuildStep

class buildCacheConfiguration(environment: String, vcsRootId: String) {
    val environment = environment
    val vcsRootId = vcsRootId

    fun buildConfiguration(providerName: String): BuildType {
        return BuildType {
            id(uniqueID(providerName))

            name = "Cache Build Dependencies"

            vcs {
                root(rootId = AbsoluteId(vcsRootId))
                cleanCheckout = true
            }

            steps {
                ConfigureGoEnv()
                step(ScriptBuildStep {
                    name = "Compile Test Binary"
                    scriptContent = """
                        mkdir -p %env.GOCACHE%
                        mkdir -p %env.GOMODCACHE%
                        go test -c -o test-binary
                    """.trimIndent()
                })
            }

            triggers {
                RunNightly(
                    nightlyTestsEnabled = true,
                    startHour = 23,
                    daysOfWeek = "*",
                    daysOfMonth = "*"
                )
            }

            failureConditions {
                errorMessage = true
                executionTimeoutMin = 60
            }

            features {
                feature(BuildCacheFeature {
                    name = "terraform-provider-azurerm-build-cache"
                    publish = true
                    use = true
                    rules = """
                        %env.GOCACHE%
                        %env.GOMODCACHE%
                    """.trimIndent()
                })
            }

            cleanup {
                baseRule {
                    artifacts(days = 7, artifactPatterns = "+:**/*")
                }
            }

            params {
                GoCache()
                ReadOnlySettings()
            }
        }
    }

    fun uniqueID(provider: String): String {
        return "%s_CACHE_%s".format(provider.uppercase(), environment.uppercase())
    }
}
