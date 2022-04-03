rm -rf /etc/docker/certs.d/{{ $.imageRepoDomainName }}
mkdir -p /etc/docker/certs.d/{{ $.imageRepoDomainName }}
export INGRESS_SECRET_NAME=$(kubectl -n {{ $.imageRepo.internal.namespace }} get secrets | grep "harbor-ingress" | awk '{print $1}')
kubectl -n {{ $.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.ca\.crt }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepoDomainName }}/ca.crt
kubectl -n {{ $.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.tls\.crt }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepoDomainName }}/{{ $.imageRepoDomainName }}.cert
kubectl -n {{ $.imageRepo.internal.namespace }} get secrets ${INGRESS_SECRET_NAME} -o jsonpath='{ .data.tls\.key }' | base64 -d > /etc/docker/certs.d/{{ $.imageRepoDomainName }}/{{ $.imageRepoDomainName }}.key
