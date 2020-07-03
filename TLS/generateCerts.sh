#!/bin/bash

# Output files
# ca.key: Certificate Authority private key file (this shouldn't be shared in real-life)
# ca.crt: Certificate Authority trust certificate (this should be shared with users in real-life)
# server.key: Server private key, password protected (this shouldn't be shared)
# server.csr: Server certificate signing request (this should be shared with the CA owner)
# server.crt: Server certificate signed by the CA (this would be sent back by the CA owner) - keep on server
# server.pem: Conversion of server.key into a format gRPC likes (this shouldn't be shared)

# Summary 
# Private files: ca.key, server.key, server.pem, server.crt
# "Share" files: ca.crt (needed by the client), server.csr (needed by the CA)

# Changes these CN's to match your hosts in your environment if needed.
SERVER_CN=localhost

echo " -------------------- Generating Certificates -------------------------"
echo ""
cd TLS

# Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
# Generate Private Key
openssl genrsa -passout pass:saanvijay -des3 -out Exampleca.key 4096
sleep 10
# Get trust certificate using private key with validity 
openssl req -passin pass:saanvijay -new -x509 -days 365 -key Exampleca.key -out Exampleca.crt -subj "/CN=${SERVER_CN}"


# Step 2: Generate the Server Private Key (scmserver.key)
openssl genrsa -passout pass:saanvijay -des3 -out scmserver.key 4096
sleep 10
# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:saanvijay -new -key scmserver.key -out scmserver.csr -subj "/CN=${SERVER_CN}"

# Step 4: Sign the certificate with the CA we created (it's called self signing) - server.crt
openssl x509 -req -passin pass:saanvijay -days 365 -in scmserver.csr -CA Exampleca.crt -CAkey Exampleca.key -set_serial 01 -out scmserver.crt 

# Step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:saanvijay -in scmserver.key -out scmserver.pem
echo ""
echo " =========================== DONE ================================"