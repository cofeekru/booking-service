pipeline {
    agent any
    stages {
        stage('Cleanup') {
            steps {
                cleanWs()
            }
        }
        stage('Checkout Code') {
            steps {
                git url: 'https://your-repo-url.git', branch: 'main'
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
