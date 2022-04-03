# 把harbor服务器({{ $.imageRepoIp }})上的证书复制到本节点的 /etc/docker/certs.d/{{ $.imageRepoDomainName }} 目录
# 证书文件包括: ca.crt, {{ $.imageRepoDomainName }}.cert, {{ $.imageRepoDomainName }}.key
# 完成harbor证书复制后，请输入 [YES] 继续安装，输入 [NO] 取消安装
