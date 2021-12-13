rm -rf /etc/docker/certs.d/{{ $.imageRepo.domainName }}
mkdir -p /etc/docker/certs.d/{{ $.imageRepo.domainName }}
export INGRESS_SECRET_NAME=$(kubectl -n {{ $.imageRepo.namespace }} get secrets | grep "harbor-ingress" | awk '{print $1}')
kubectl -n {{ $.imageRepo.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.ca\.crt }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepo.domainName }}/ca.crt
kubectl -n {{ $.imageRepo.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.tls\.crt }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepo.domainName }}/{{ $.imageRepo.domainName }}.cert
kubectl -n {{ $.imageRepo.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.tls\.key }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepo.domainName }}/{{ $.imageRepo.domainName }}.key
