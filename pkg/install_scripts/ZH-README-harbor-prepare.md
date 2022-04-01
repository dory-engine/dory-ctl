# 把harbor服务器({{ $.imageRepo.external.ip }})上的证书复制到本节点的 /etc/docker/certs.d/{{ $.imageRepo.external.url }} 目录
# 证书文件包括: ca.crt, {{ $.imageRepo.external.url }}.cert, {{ $.imageRepo.external.url }}.key

# 完成harbor证书复制后，请输入 [YES] 继续安装，输入 [NO] 取消安装
