# copy harbor server ({{ $.imageRepoIp }}) certificates to this node /etc/docker/certs.d/{{ $.imageRepoDomainName }} directory
# certificates are: ca.crt, {{ $.imageRepoDomainName }}.cert, {{ $.imageRepoDomainName }}.key
# after finish harbor certificates copy, please input [YES] to go on, input [NO] to cancel
