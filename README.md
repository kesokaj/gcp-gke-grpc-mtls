## CA + Server certificate
````
openssl genrsa -out ca.key 2048
openssl req -new -x509 -key ca.key -out ca.pem -days 3650 -subj "/CN=My Example CA"
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=example.com"

cat <<EOF > v3.ext
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = example.com
DNS.2 = grpc.example.com
EOF

openssl x509 -req -in server.csr -out server.crt -CA ca.pem -CAkey ca.key -CAcreateserial -days 365 -sha256 -extfile v3.ext
````

## Client certificate
````
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "/CN=client.example.com"
openssl x509 -req -in client.csr -out client.crt -CA ca.pem -CAkey ca.key -CAcreateserial -days 365 -sha256

````

## Create kubernetes secrets
````
kubectl create secret tls grpc-crt --cert=server.crt --key=server.key
kubectl create secret generic grpc-ca --from-file=ca.pem
````

## Deploy the apps for internal use
````
kubectl apply -f manifests/
````

## Run your tests with grpcurl
````
### mtls ###
grpcurl -cacert ca.pem -cert client.crt -key client.key -use-reflection -d '{"name": "World"}' -H 'Content-Type: application/json' grpc.example.com:443 helloworld.Greeter.SayHello

### Insecure ###
grpcurl -plaintext -d '{"name": "World"}' -H 'Content-Type: application/json' grpc.example.com:80 helloworld.Greeter.SayHello

````