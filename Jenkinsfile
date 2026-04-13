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
                withCredentials([file(credentialsId: 'kubeconfig-minikube', variable: 'KUBE_CONFIG')]) {
                    sh '''
                    mkdir -p ~/.kube
                    cp $KUBE_CONFIG ~/.kube/config
                    chmod 600 ~/.kube/config
                    ./deploy/apply.sh
                '''
                }
        }
    }
}
