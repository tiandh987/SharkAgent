#!/usr/bin/env bash

# BASH_SOURCE[0] 获取当前脚本名称
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

# source common.sh
source "${IAM_ROOT}/scripts/common.sh"

# LOCAL_OUTPUT_ROOT 定义在 common.sh 中
readonly LOCAL_OUTPUT_CONFIGPATH="${LOCAL_OUTPUT_ROOT}/configs"
mkdir -p ${LOCAL_OUTPUT_CONFIGPATH}

cd ${IAM_ROOT}/scripts

export IAM_APISERVER_INSECURE_BIND_ADDRESS=0.0.0.0
export IAM_AUTHZ_SERVER_INSECURE_BIND_ADDRESS=0.0.0.0

# 集群内通过 kubernetes 服务名访问
export IAM_APISERVER_HOST=iam-apiserver
export IAM_AUTHZ_SERVER_HOST=iam-authz-server
export IAM_PUMP_HOST=iam-pump
export IAM_WATCHER_HOST=iam-watcher

# 配置 CA 证书路径
export CONFIG_USER_CLIENT_CERTIFICATE=/etc/iam/cert/admin.pem
export CONFIG_USER_CLIENT_KEY=/etc/iam/cert/admin-key.pem

export CONFIG_SERVER_CERTIFICATE_AUTHORITY=/etc/iam/cert/ca.pem

# 生成配置
for comp in iam-apiserver iam-authz-server iam-pump iam-watcher iamctl
do
  iam::log::info "generate ${LOCAL_OUTPUT_CONFIGPATH}/${comp}.yaml"
  ./genconfig.sh install/environment.sh ../configs/${comp}.yaml > ${LOCAL_OUTPUT_CONFIGPATH}/${comp}.yaml
done
