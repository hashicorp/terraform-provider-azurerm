@Library(['pipeline-toolbox', 'iac-pipeline-shared-lib']) _

pipeline {
    agent any

    stages {        
        stage('Test Build') {
            steps {
              script {
                docker.image("dockerhub.rnd.amadeus.net/docker-production/iac/terraform-automation-azr:2.7.2").inside("-u iacuser") {
                  withCredentials([
                    string(credentialsId: 'GPG_FINGERPRINT', variable: 'GPG_FINGERPRINT'),
                    file(credentialsId: 'ash-gpg-key', variable: 'ASH_GPG_KEY'),
                    usernamePassword(credentialsId: '	goreleaser-artifactory-creds', usernameVariable: 'ARTIFACTORY_PRODUCTION_USERNAME', passwordVariable: 'ARTIFACTORY_PRODUCTION_SECRET')
                  ]) {
                    sh 'cp "${ASH_GPG_KEY}" .'
                    sh 'gpg --import ash-gpg-key && rm -rf ash-gpg-key'
                    sh 'wget -q -O /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/v1.11.4/goreleaser_Linux_x86_64.tar.gz'
                    sh 'tar -xf /tmp/goreleaser.tar.gz --directory /tmp/'
                    sh '/tmp/goreleaser release --rm-dist --parallelism=2'
                    sh 'ls -lart dist/'
                  }                  
                }           
              }
            }                        
        }
    }
    
}