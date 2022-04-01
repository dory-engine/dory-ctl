rm -rf /etc/docker/certs.d/{{ $.imageRepo.internal.domainName }}
mkdir -p /etc/docker/certs.d/{{ $.imageRepo.internal.domainName }}
export INGRESS_SECRET_NAME=$(kubectl -n {{ $.imageRepo.internal.namespace }} get secrets | grep "harbor-ingress" | awk '{print $1}')
kubectl -n {{ $.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.ca\.crt }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepo.internal.domainName }}/ca.crt
kubectl -n {{ $.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.tls\.crt }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepo.internal.domainName }}/{{ $.imageRepo.internal.domainName }}.cert
kubectl -n {{ $.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.tls\.key }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepo.internal.domainName }}/{{ $.imageRepo.internal.domainName }}.key
