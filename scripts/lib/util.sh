#!/usr/bin/env bash

iam::util::host_arch() {
  local host_arch
  case "$(uname -m)" in
    x86_64*)
      host_arch=amd64
      ;;
    i?86_64*)
      host_arch=amd64
      ;;
    amd64*)
      host_arch=amd64
      ;;
    aarch64*)
      host_arch=arm64
      ;;
    arm64*)
      host_arch=arm64
      ;;
    arm*)
      host_arch=arm
      ;;
    i?86*)
      host_arch=x86
      ;;
    s390x*)
      host_arch=s390x
      ;;
    ppc64le*)
      host_arch=ppc64le
      ;;
    *)
      iam::log::error "Unsupported host arch. Must be x86_64, 386, arm, arm64, s390x or ppc64le."
      exit 1
      ;;
  esac
  echo "${host_arch}"
}

# 如果 PATH 中不存在 cfssl/cfssljson/cfssl-certinfo，则将它们下载到 $1 目录中
#
# 假设变量：
# $1 (cfssl 目录) (可选)
#
# Sets：
#   CFSSL_BIN: 已安装 cfssl 二进制文件的路径
#   CFSSLJSON_BIN：安装的cfssljson二进制文件的路径
#   CFSSLCERTINFO_BIN: 安装的 cfssl-certinfo 二进制文件的路径
#
# cfssl：证书签发工具。
# cfssljson：将 cfssl 生成的证书（json 格式）变为文件承载式证书。
function iam::util::ensure-cfssl() {
  # [tian@dev workspace]$ command -v cfssl
  # /home/tian/bin/cfssl
  #
  if command -v cfssl &>/dev/null && command -v cfssljson &>/dev/null && command -v cfssl-certinfo &>/dev/null; then
    CFSSL_BIN=$(command -v cfssl)
    CFSSLJSON_BIN=$(command -v cfssljson)
    CFSSLCERTINFO_BIN=$(command -v cfssl-certinfo)
    return 0
  fi

  host_arch=$(iam::util::host_arch)

  if [[ "${host_arch}" != "amd64" ]]; then
    echo "Cannot download cfssl on non-amd64 hosts and cfssl does not appear to be installed."
    echo "Please install cfssl, cfssljson and cfssl-certinfo and verify they are in \$PATH."
    echo "Hint: export PATH=\$PATH:\$GOPATH/bin; go get -u github.com/cloudflare/cfssl/cmd/..."
    exit 1
  fi

  # Create a temp dir for cfssl if no directory was given
  local cfssldir=${1:-}
  if [[ -z "${cfssldir}" ]]; then
    cfssljson="$HOME/bin"
  fi

  mkdir -p "${cfssldir}"
  pushd "${cfssldir}" > /dev/null || return 1

  echo "Unable to successfully run 'cfssl' from ${PATH}; downloading instead..."
  kernel=$(uname -s)
  case "${kernel}" in
    Linux)
      curl --retry 10 -L -o cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
      curl --retry 10 -L -o cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
      curl --retry 10 -L -o cfssl-certinfo https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64
      ;;
    Darwin)
      curl --retry 10 -L -o cfssl https://pkg.cfssl.org/R1.2/cfssl_darwin-amd64
      curl --retry 10 -L -o cfssljson https://pkg.cfssl.org/R1.2/cfssljson_darwin-amd64
      curl --retry 10 -L -o cfssl-certinfo https://pkg.cfssl.org/R1.2/cfssl-certinfo_darwin-amd64
      ;;
    *)
      echo "Unknown, unsupported platform: ${kernel}." >&2
      echo "Supported platforms: Linux, Darwin." >&2
      exit 2
  esac

  chmod +x cfssl || true
  chmod +x cfssljson || true
  chmod +x cfssl-certinfo || true

  CFSSL_BIN="${cfssldir}/cfssl"
  CFSSLJSON_BIN="${cfssldir}/cfssljson"
  CFSSLCERTINFO_BIN="${cfssldir}/cfssl-certinfo"

  # [ -x FILE ] 如果 FILE 存在且是可执行的则为真。
  if [[ ! -x ${CFSSL_BIN} || ! -x ${CFSSLJSON_BIN} || ! -x ${CFSSLCERTINFO_BIN} ]]; then
    echo "Failed to download 'cfssl'."
    echo "Please install cfssl, cfssljson and cfssl-certinfo and verify they are in \$PATH."
    echo "Hint: export PATH=\$PATH:\$GOPATH/bin; go get -u github.com/cloudflare/cfssl/cmd/..."
    exit 1
  fi
  popd > /dev/null || return 1
}