pipeline {
    agent any
    stages {
        stage('Checkout Code') {
            steps {
                git url: 'http://github.com/cofeekru/booking-service.git', branch: 'master'
            }
        }

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

        stage('Deploy to Kubernetes') {
            steps {
                sh '''
                    kubectl apply -f deploy/postgres/postgres-pvc.yaml
                    kubectl apply -f deploy/postgres/postgres-deployment.yaml
                    kubectl apply -f deploy/postgres/postgres-service.yaml

                    kubectl apply -f deploy/app/deployment.yaml
                    kubectl apply -f deploy/app/service.yaml
                '''
            }
        }
        stage('Cleanup') {
            steps {
                cleanWs()
            }
        }  

    }
}
