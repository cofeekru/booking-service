pipeline {
    agent any
    stages {
        stage('Cleanup') {
            steps {
                cleanWs()
            }
        }
        stage('Verify Workspace') {
            steps {
                sh 'pwd'
                sh 'ls -la'
                sh 'find . -name "Dockerfile*" -o -name "dockerfile*"'
            }
        }
        
        stage('Build') {
            steps {
                sh 'docker build -t my-app:latest' 
            }
        }

        stage('Push to docker registry') {
            steps {
                sh 'docker push localhost:5000/my-app:latest' 
            }
        }

    }
}
