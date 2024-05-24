#!/bin/bash

echo "Lets create some certs!"


mkdir certs
# Root CA

#generate ca in /certs
cfssl gencert -initca ./deployment/ca-csr.json | cfssljson -bare ./certs/ca

#generate certificate in 

cfssl gencert \
  -ca=./certs/ca.pem \
  -ca-key=./certs/ca-key.pem \
  -config=./deployment/ca-config.json \
  -hostname="harrison-admission.default.svc" \
  -profile=default \
  ./deployment/ca-csr.json | cfssljson -bare ./certs/harrison-admission


echo "Lets create the secret for the deployment to use. The Go server will load these from /etc/certs."

echo "Using the default namespace for this. In production you prob will want to use another one."

kubectl delete secret admission
kubectl create secret tls admission \
        --cert "./certs/harrison-admission.pem" \
        --key "./certs/harrison-admission-key.pem"

echo "Installing webhooks"

CA_BUNDLE=$(cat certs/ca.pem | base64 | tr -d '\n')
sed -e 's@${CA_BUNDLE}@'"$CA_BUNDLE"'@g' <"deployment/webhooks.yaml" | kubectl apply -f -