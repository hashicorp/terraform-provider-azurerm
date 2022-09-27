@Library(['pipeline-toolbox', 'iac-pipeline-shared-lib']) _

pipeline {
    agent any

    stages {        
        stage('Test Build') {
            script {
              docker.image("dockerhub.rnd.amadeus.net/docker-production-iac/iac/terraform-automation-azr:master-004346b").inside("-u iacuser") {
                sh 'wget -O /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/v1.11.4/goreleaser_Linux_x86_64.tar.gz'
                sh 'tar -xf -C /tmp/ /tmp/goreleaser.tar.gz'
                sh '/tmp/goreleaser -h'
              }           
            }            
        }
    }
    
}