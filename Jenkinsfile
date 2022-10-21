/* groovylint-disable DuplicateMapLiteral, LineLength */
/* groovylint-disable-next-line CompileStatic, NoDef, UnusedVariable, VariableName, VariableTypeRequired */
@Library(['pipeline-toolbox', 'iac-pipeline-shared-lib']) _

node {
    try {
        stage('Setup') {
            artifactName = env.BITBUCKET_REPOSITORY
            checkoutGit()
            terraformAutoImageVersion = '1.0.2'
            baseTerraformAutoImage = "docker-production-iac/iac/tf-plugin-builder:${terraformAutoImageVersion}"
            baseVersion = 'v1.0'
            newVersion = newVersion(baseVersion) + '-amadeus'
            registry = 'repository.adp.amadeus.net'
            baseImage = 'maven:3.6.3-jdk-11'
        // pageId for https://rndwww.nce.amadeus.net/confluence/display/IBSDC/IaC+Release+Notes
        //releaseNotesOptions = ['spaceKey': 'IBSDC', 'parentPageId': 1654937615]
        }

        stage('QA tests and goreleaser to release locally (not pushing to artifactory yet)') {
            /* groovylint-disable-next-line UnnecessaryGetter */
            when(isPullRequest()) {
                docker.withRegistry("https://${registry}") {
                    /* groovylint-disable-next-line NestedBlockDepth */
                    docker.image(baseTerraformAutoImage).inside {
                        withCredentials([
                          usernamePassword(credentialsId: 'goreleaser-artifactory-creds', usernameVariable: 'ARTIFACTORY_PRODUCTION_USERNAME', passwordVariable: 'ARTIFACTORY_PRODUCTION_SECRET'),
                          usernamePassword(credentialsId: 'RND-ARTIFACTORY-TOKEN', usernameVariable: 'RND_ARTIFACTORY_USER', passwordVariable: 'RND_ARTIFACTORY_TOKEN'),
                          usernamePassword(credentialsId: 'MUC_ARTIFACTORY_REGISTRY_TOKEN', usernameVariable: 'MUC_REGISTRY_HOST', passwordVariable: 'MUC_ARTIFACTORY_TOKEN')
                        /* groovylint-disable-next-line NestedBlockDepth */
                        ]) {
                            sh '''
                              make fmtcheck
                              make test
                              #not running acceptance tests now
                              #make testacc
                              echo -e "credentials \\"$MUC_REGISTRY_HOST\\" {\n   token = \\"$MUC_ARTIFACTORY_TOKEN\\"\n}\n" > .terraformrc
                              goreleaser release --snapshot --rm-dist --config .goreleaser-jenkins.yml --parallelism=2
                        '''
                        }
                    }
                }
            }
        }

        stage('Increment Tag && gorelease to artifactory') {
            withCredentials([
                /* groovylint-disable-next-line DuplicateStringLiteral */
                usernamePassword(credentialsId: 'goreleaser-artifactory-creds', usernameVariable: 'ARTIFACTORY_PRODUCTION_USERNAME', passwordVariable: 'ARTIFACTORY_PRODUCTION_SECRET'),
                /* groovylint-disable-next-line DuplicateStringLiteral */
                usernamePassword(credentialsId: 'RND-ARTIFACTORY-TOKEN', usernameVariable: 'RND_ARTIFACTORY_USER', passwordVariable: 'RND_ARTIFACTORY_TOKEN'),
                /* groovylint-disable-next-line DuplicateStringLiteral */
                usernamePassword(credentialsId: 'MUC_ARTIFACTORY_REGISTRY_TOKEN', usernameVariable: 'MUC_REGISTRY_HOST', passwordVariable: 'MUC_ARTIFACTORY_TOKEN')
          ]) {
                /* groovylint-disable-next-line UnnecessaryGetter */
                when(env.BRANCH_NAME == 'master' || isReleasedBranch()) {
                    pushNewVersionTag(newVersion, baseVersion, releaseNotesOptions)
                    /* groovylint-disable-next-line DuplicateStringLiteral, NestedBlockDepth */
                    docker.withRegistry("https://${registry}") {
                        /* groovylint-disable-next-line NestedBlockDepth */
                        docker.image(baseTerraformAutoImage).inside {
                            /* groovylint-disable-next-line GStringExpressionWithinString */
                            sh '''
                              echo -e "credentials \\"$MUC_REGISTRY_HOST\\" {\n   token = \\"$MUC_ARTIFACTORY_TOKEN\\"\n}\n" > .terraformrc
                              goreleaser release --rm-dist --config .goreleaser-jenkins.yml --parallelism=2
                            '''
                        }
                    }
                }
          }
        }
    } catch (err) {
        echo "Caught: ${err}"
        currentBuild.result = 'FAILURE'
    } finally {
        echo 'Done'
    }
}
