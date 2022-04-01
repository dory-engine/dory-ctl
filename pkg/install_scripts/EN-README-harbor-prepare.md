# copy harbor server ({{ $.imageRepo.external.ip }}) certificates to this node /etc/docker/certs.d/{{ $.imageRepo.external.url }} directory
# certificates are: ca.crt, {{ $.imageRepo.external.url }}.cert, {{ $.imageRepo.external.url }}.key
# after finish harbor certificates copy, please input [YES] to go on, input [NO] to cancel
