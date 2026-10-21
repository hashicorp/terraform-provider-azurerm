import jetbrains.buildServer.configs.kotlin.AbsoluteId
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.buildFeatures.BuildCacheFeature
import jetbrains.buildServer.configs.kotlin.buildSteps.ScriptBuildStep

// The test builds no longer fetch what this publishes: each agent keeps its own Go caches instead, see
// GoCache(). It's left running so that fetching can be switched back on if that doesn't work out, and
// can be removed once it's clear that it does. Meanwhile it warms the cache of whichever agent runs it.
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
                    startHour = 22,
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
                    use = false
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
                // ConfigureGoEnv() only runs when this is set, and without it the agent's default Go is used
                text("env.SCHEDULE_MATCHES", "true")
            }
        }
    }

    fun uniqueID(provider: String): String {
        return "%s_CACHE_%s".format(provider.uppercase(), environment.uppercase())
    }
}
