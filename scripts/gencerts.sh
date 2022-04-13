#!/usr/bin/env bash

# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${IAM_ROOT}/scripts/lib/init.sh"

# OUT_DIR can come in from the Makefile, so honor it.
readonly LOCAL_OUTPUT_ROOT="${IAM_ROOT}/${OUT_DIR:-_output}"
readonly LOCAL_OUTPUT_CAPATH="${LOCAL_OUTPUT_ROOT}/cert"

# Hostname for the cert
readonly CERT_HOSTNAME="${CERT_HOSTNAME:-iam.api.marmotedu.com,iam.authz.marmotedu.com},127.0.0.1,localhost"

# 运行 cfssl命令生成 iam 服务的证书文件，证书文件会保存在 $1目录下。
#
# 参数：
#   $1（证书文件保存的目录）
#   $2（证书文件名的前缀）
function generate-iam-cert() {
  # local一般用于局部变量声明，多在在函数内部使用。
  #
  #（1）shell脚本中定义的变量是global的，其作用域从被定义的地方开始，到shell结束或被显示删除的地方为止。
  #
  #（2）函数定义的变量可以被显示定义成local的，其作用域局限于函数内。
  #     但请注意，函数的参数是local的。
  #
  #（3）如果同名，Shell函数定义的local变量会屏蔽脚本定义的global变量。
  #
  local cert_dir=${1}
  local prefix=${2:-}

  mkdir -p "${cert_dir}"
  pushd "${cert_dir}"

  iam::util::ensure-cfssl

  # -r file 用户可读为真
  #
  # CA 配置文件是用来配置根证书的使用场景 (profile) 和具体参数 (usage、过期时间、服务端认证、客户端认证、加密等)，
  # 可以在签名其它证书时用来指定特定场景。
  #
  # JSON 配置中，有一些字段解释如下：
  #   signing：表示该证书可用于签名其它证书（生成的 ca.pem 证书中 CA=TRUE）。
  #   server auth：表示 client 可以用该证书对 server 提供的证书进行验证。
  #   client auth：表示 server 可以用该证书对 client 提供的证书进行验证。
  #   expiry：876000h，证书有效期设置为 100 年。
  #
  if [ ! -r "ca-config.json" ]; then
  # 创建 CA 配置文件
    cat >ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "iam": {
        "usages": [
          "signing",
          "key encipherment",
          "server auth",
          "client auth"
        ],
        "expiry": "876000h"
      }
    }
  }
}
EOF
  fi

  # 创建证书签名请求文件。
  # JSON 配置中，有一些字段解释如下：
  #   CN：Common Name，iam-apiserver 从证书中提取该字段作为请求的用户名 (User Name) ，浏览器使用该字段验证网站是否合法。
  #   C：Country，国家。
  #   ST：State，省份。
  #   L：Locality (L) or City，城市。
  #   O：Organization，iam-apiserver 从证书中提取该字段作为请求用户所属的组 (Group)。
  #   OU：Company division (or Organization Unit – OU)，部门 / 单位。
  #
  # 除此之外，还有两点需要我们注意：
  #   不同证书 csr 文件的 CN、C、ST、L、O、OU 组合必须不同，否则可能出现 PEER'S CERTIFICATE HAS AN INVALID SIGNATURE 错误。
  #   后续创建证书的 csr 文件时，CN、OU 都不相同（C、ST、L、O 相同），以达到区分的目的。
  #
  if [ ! -r "ca-csr.json" ]; then
    cat >ca-csr.json <<EOF
{
  "CN": "iam-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "marmotedu",
      "OU": "iam"
    }
  ],
  "ca": {
    "expiry": "876000h"
  }
}
EOF
  fi

  # 创建 CA 证书 和 私钥
  if [[ ! -r "ca.pem" || ! -r "ca-key.pem" ]]; then
    # cfssl gencert 命令会创建运行 CA 所必需的文件
    #   ca-key.pem（私钥）
    #   ca.pem（证书）
    #   ca.csr（证书签名请求），用于交叉签名或重新签名。
    ${CFSSL_BIN} gencert -initca ca-csr.json | ${CFSSLJSON_BIN} -bare ca -

    # 创建完之后，我们可以通过 cfssl certinfo 命名查看 cert 和 csr 信息：
    #   $ cfssl certinfo -cert ${IAM_CONFIG_DIR}/cert/ca.pem  # 查看 cert(证书信息)
    #   $ cfssl certinfo -csr ${IAM_CONFIG_DIR}/cert/ca.csr   # 查看 CSR(证书签名请求)信息
  fi

  if [[ -z "${prefix}" ]];then
    return 0
  fi

  # 创建 iam-apiserver、iam-auth-server、admin 证书和私钥。
  # hosts 字段是用来指定授权使用该证书的 IP 和域名列表。
  # -hostname=""：证书的hostname，可以是逗号分隔的hostname列表;
  # -ca="": 用于签署新证书的 CA;
  # -ca-key=""：CA 私钥;
  # -config=""：配置文件的路径;
  # -profile=""：要使用的签名配置文件;
  #
  # 关于减号【-】的作用：
  #   在管道命令当中，常常会使用前一个命令的 stdout 作为这次的 stdin，
  #   该 stdin 与 stdout 可以利用减号 “-” 来替代。
  echo "Generate "${prefix}" certificates..."
  echo '{"CN":"'"${prefix}"'","hosts":[],"key":{"algo":"rsa","size":2048},"names":[{"C":"CN","ST":"BeiJing","L":"BeiJing","O":"marmotedu","OU":"'"${prefix}"'"}]}' \
    | ${CFSSL_BIN} gencert -hostname="${CERT_HOSTNAME},${prefix}" -ca=ca.pem -ca-key=ca-key.pem \
    -config=ca-config.json -profile=iam - | ${CFSSLJSON_BIN} -bare "${prefix}"

  # the popd will access `directory stack`, no `real` parameters is actually needed
  # shellcheck disable=SC2119
  popd
}

# 为 iam 组件生成 SSL 证书。 使用 cfssl 程序。
#
# Assumed vars:
#   IAM_TEMP：临时目录
#
# Args:
#   $1（证书文件名的前缀）
#
# 如果 CA cert/key 为空，该函数也会为 CA 生成证书。
#
# 变量集：
#   IAM_CA_KEY_BASE64
#   IAM_CA_CERT_BASE64
#   IAM_APISERVER_KEY_BASE64
#   IAM_APISERVER_CERT_BASE64
#   IAM_AUTHZ_SERVER_KEY_BASE64
#   IAM_AUTHZ_SERVER_CERT_BASE64
#   IAM_ADMIN_KEY_BASE64
#   IAM_ADMIN_CERT_BASE64
#
function create-iam-certs {
  local prefix=${1}

  iam::util::ensure-temp-dir

  generate-iam-cert "${IAM_TEMP}/cfssl" ${prefix}

  pushd "${IAM_TEMP}/cfssl"
  IAM_CA_KEY_BASE64=$(cat "ca-key.pem" | base64 | tr -d '\r\n')
  IAM_CA_CERT_BASE64=$(cat "ca.pem" | gzip | base64 | tr -d '\r\n')

  case "${prefix}" in
    iam-apiserver)
      IAM_APISERVER_KEY_BASE64=$(cat "iam-apiserver-key.pem" | base64 | tr -d '\r\n')
    	IAM_APISERVER_CERT_BASE64=$(cat "iam-apiserver.pem" | gzip | base64 | tr -d '\r\n')
    	;;
    iam-authz-server)
			IAM_AUTHZ_SERVER_KEY_BASE64=$(cat "iam-authz-server-key.pem" | base64 | tr -d '\r\n')
			IAM_AUTHZ_SERVER_CERT_BASE64=$(cat "iam-authz-server.pem" | gzip | base64 | tr -d '\r\n')
			;;
    admin)
			IAM_ADMIN_KEY_BASE64=$(cat "admin-key.pem" | base64 | tr -d '\r\n')
			IAM_ADMIN_CERT_BASE64=$(cat "admin.pem" | gzip | base64 | tr -d '\r\n')
			;;
    *)
      echo "Unknow, unsupported iam certs type:: ${prefix}" >&2
      echo "Supported type: iam-apiserver, iam-authz-server, admin" >&2
			exit 2
  esac

  popd
}

# $* 传递给脚本或函数的所有参数
$*