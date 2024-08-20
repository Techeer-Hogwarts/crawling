pipeline {
    agent any
    
    environment {
        repository = "suhach0523/techeerism-crawler"  //docker hub id와 repository 이름
        DOCKERHUB_CREDENTIALS = credentials('docker-hub') // jenkins에 등록해 놓은 docker hub credentials 이름
        IMAGE_TAG = ""
    }

    stages {
        stage('Checkout') {
            steps {
                cleanWs()
                git branch: 'main', url: "https://github.com/Techeer-Hogwarts/crawling.git"
            }
        }

        stage('Test') {
            steps {
                script {
                    // Step 0: Check Go installation
                    echo 'Checking Go installation...'
                    sh 'go version'
                    // Step 1: Run Go linting
                    echo 'Running Go linting...'
                    sh 'go clean -modcache'
                    sh 'golangci-lint run --timeout 5m --verbose'

                    // Step 2: Run Go unit tests
                    echo 'Running Go unit tests...'
                    sh 'go test ./...'

                    // Step 3: Check Go code formatting
                    echo 'Checking Go code formatting...'
                    sh 'go fmt ./...'

                    // Step 4: Check Docker installation
                    echo 'Checking Docker installation...'
                    sh 'docker --version'
                }
            }
        }


        stage('Set Image Tag') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'main') {
                        IMAGE_TAG = "1.0.${BUILD_NUMBER}"
                    } else {
                        IMAGE_TAG = "0.0.${BUILD_NUMBER}"
                    }
                    echo "Image tag set to: ${IMAGE_TAG}"
                }
            }
        }

        stage('Building our image') { 
            steps { 
                script { 
                    sh "docker build -t ${repository}:${IMAGE_TAG} ." // docker build
                }
                slackSend message: "Build Started - ${env.JOB_NAME} ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
            } 
        }
        
        stage('Login') {
            steps {
                sh "echo $DOCKERHUB_CREDENTIALS_PSW | docker login -u $DOCKERHUB_CREDENTIALS_USR --password-stdin"
            }
        }

        stage('Image Push') { 
            steps { 
                script {
                    sh "docker push ${repository}:${IMAGE_TAG}"
                } 
            }
        } 

        stage('Clean up') { 
            steps { 
                sh "docker rmi ${repository}:${IMAGE_TAG}" // docker image 제거
            }
        } 
    }

    post {
        always {
            cleanWs(cleanWhenNotBuilt: false,
                    deleteDirs: true,
                    disableDeferredWipeout: true,
                    notFailBuild: true,
                    patterns: [[pattern: '.gitignore', type: 'INCLUDE'],
                            [pattern: '.propsfile', type: 'EXCLUDE']])
        }
        success {
            echo 'Build and deployment successful!'
            slackSend message: "Build deployed successfully - ${env.JOB_NAME} ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
        }
        failure {
            echo 'Build or deployment failed.'
            slackSend failOnError: true, message: "Build failed  - ${env.JOB_NAME} ${env.BUILD_NUMBER} (<${env.BUILD_URL}|Open>)"
        }
    }
}