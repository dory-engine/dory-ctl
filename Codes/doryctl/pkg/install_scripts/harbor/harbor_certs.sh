# 设置harbor域名
export HARBOR_CONFIG_DOMAIN_NAME={{ $.harbor.domainName }}
# 设置证书存放的相对路径
export HARBOR_CONFIG_CERT_PATH={{ $.harbor.certsDir }}

rm -rf ${HARBOR_CONFIG_CERT_PATH}/
mkdir -p ${HARBOR_CONFIG_CERT_PATH}/
cd ${HARBOR_CONFIG_CERT_PATH}/

openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -sha512 -days 3650 -subj "/CN=${HARBOR_CONFIG_DOMAIN_NAME}" -key ca.key -out ca.crt
openssl genrsa -out ${HARBOR_CONFIG_DOMAIN_NAME}.key 4096
openssl req -sha512 -new -subj "/CN=${HARBOR_CONFIG_DOMAIN_NAME}" -key ${HARBOR_CONFIG_DOMAIN_NAME}.key -out ${HARBOR_CONFIG_DOMAIN_NAME}.csr
cat > v3.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=${HARBOR_CONFIG_DOMAIN_NAME}
EOF
openssl x509 -req -sha512 -days 3650 -extfile v3.ext -CA ca.crt -CAkey ca.key -CAcreateserial -in ${HARBOR_CONFIG_DOMAIN_NAME}.csr -out ${HARBOR_CONFIG_DOMAIN_NAME}.crt
openssl x509 -inform PEM -in ${HARBOR_CONFIG_DOMAIN_NAME}.crt -out ${HARBOR_CONFIG_DOMAIN_NAME}.cert
# echo "[INFO] # check harbor certificates info"
# echo "[CMD] openssl x509 -noout -text -in ${HARBOR_CONFIG_DOMAIN_NAME}.crt"
# openssl x509 -noout -text -in ${HARBOR_CONFIG_DOMAIN_NAME}.crt

# 更新/etc/docker/certs.d/对应的harbor证书
echo "update docker harbor certificates"
rm -rf /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
mkdir -p /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ${HARBOR_CONFIG_DOMAIN_NAME}.cert /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ${HARBOR_CONFIG_DOMAIN_NAME}.key /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cp ca.crt /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
ls -al /etc/docker/certs.d/${HARBOR_CONFIG_DOMAIN_NAME}/
cd ..
