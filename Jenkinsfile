pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'docker build -t my-app:latest .' 
            }
        }

        stage('Push to docker registry') {
            steps {
                sh 'docker push localhost:5000/my-app:latest' 
            }
        }

    }
}
