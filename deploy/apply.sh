#!//bin/bash

kubectl apply -f deploy/postgres/postgres-pvc.yaml
kubectl apply -f deploy/postgres/postgres-deployment.yaml
kubectl apply -f deploy/postgres/postgres-service.yaml

kubectl apply -f deploy/app/deployment.yaml
kubectl apply -f deploy/app/service.yaml
