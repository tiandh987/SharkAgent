#!/usr/bin/env bash

# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

source "${IAM_ROOT}/scripts/install/common.sh"

# =================新申请服务器配置 & Go开发环境配置====================

# 初始化新申请的Linux服务器，使其成为一个友好的开发机
function iam::install::init_into_go_env()
{
  # 1. Linux 服务器基本配置
  iam::install::prepare_linux || return 1

  # 2. Go 编译环境安装和配置
  iam::install::go || return 1

  # 3. Go 开发 IDE 安装和配置
  #iam::install::vim_ide || return 1

  iam::log::info "initialize linux to go development machine  successfully"
}

# 申请服务器，登录 going 用户后，配置 $HOME/.bashrc 文件
iam::install::prepare_linux()
{
  if [[ -f $HOME/.bashrc ]];then
    cp $HOME/.bashrc $HOME/bashrc.iam.backup
  fi

  # 1. 配置 $HOME/.bashrc
  cat << 'EOF' > $HOME/.bashrc
# .bashrc

# User specific aliases and functions

alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

# Source global definitions
if [ -f /etc/bashrc ]; then
    . /etc/bashrc
fi

if [ ! -d $HOME/workspace ]; then
    mkdir -p $HOME/workspace
fi

# User specific environment
# Basic envs
export LANG="en_US.UTF-8"          # 设置系统语言为 en_US.UTF-8，避免终端出现中文乱码
export PS1='[\u@dev \W]\$ '        # 默认的 PS1 设置会展示全部的路径，为了防止过长，这里只展示："用户名@dev 最后的目录名"
export WORKSPACE="$HOME/workspace" # 设置工作目录
export PATH=$HOME/bin:$PATH        # 将 $HOME/bin 目录加入到 PATH 变量中

# Default entry folder
cd $WORKSPACE # 登录系统，默认进入 workspace 目录

# User specific aliases and functions
EOF

  # 创建工作目录
  mkdir -p $HOME/workspace

  # 3. 安装依赖包
  iam::common::sudo "yum -y install make autoconf automake cmake perl-CPAN libcurl-devel libtool gcc gcc-c++ glibc-headers zlib-devel git-lfs telnet ctags lrzsz jq"

  # 4. 安装 Git
  cd /tmp
  wget https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.30.2.tar.gz
  tar -xvzf git-2.30.2.tar.gz
  cd git-2.30.2/
  ./configure
  make
  iam::common::sudo "make install"

  # 5. 配置Git
  cat << 'EOF' >> $HOME/.bashrc
# Configure for git
export PATH=/usr/local/libexec/git-core:$PATH
EOF

  git config --global user.name "Lingfei Kong"           # 用户名改成自己的
  git config --global user.email "colin404@foxmail.com"  # 邮箱改成自己的
  git config --global credential.helper store            # 设置 Git，保存用户名和密码
  git config --global core.longpaths true                # 解决 Git 中 'Filename too long' 的错误
  git config --global core.quotepath off
  #git config --global url."https://github.com.cnpmjs.org/".insteadOf "https://github.com/"
  git lfs install --skip-repo

  source $HOME/.bashrc
  iam::log::info "prepare linux basic environment successfully"
}

# 安装 golang
function iam::install::go()
{
  iam::install::go_command || return 1
  iam::install::protobuf || return 1

  iam::log::info "install go develop environment successfully"
}

# Go 编译环境安装和配置
function iam::install::go_command()
{
  # 检查 go 是否安装
  command -v go &>/dev/null && return 0

  # 1. 下载 go1.17.2 版本的Go安装包
  wget -P /tmp/ https://golang.google.cn/dl/go1.17.2.linux-amd64.tar.gz

  # 2. 安装Go
  mkdir -p $HOME/go
  tar -xvzf /tmp/go1.17.2.linux-amd64.tar.gz -C $HOME/go
  mv $HOME/go/go $HOME/go/go1.17.2

  # 3. 配置Go环境变量
  cat << 'EOF' >> $HOME/.bashrc
# Go envs
export GOVERSION=go1.17.2                 # Go 版本设置
export GO_INSTALL_DIR=$HOME/go            # Go 安装目录
export GOROOT=$GO_INSTALL_DIR/$GOVERSION  # GOROOT 设置
export GOPATH=$WORKSPACE/golang           # GOPATH 设置
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH # 将 Go 语言自带的和通过 go install 安装的二进制文件加入到 PATH 路径中
export GO111MODULE="on"                   # 开启 Go moudles 特性
export GOPROXY=https://mirrors.aliyun.com/goproxy,https://goproxy.cn,direct # 安装 Go 模块时，代理服务器设置
export GOPRIVATE=github.com               # 指定不走代理的 Go 包域名
export GOSUMDB=off                        # 关闭校验 Go 依赖包的哈希值
EOF
  source $HOME/.bashrc
  iam::log::info "install go compile tool successfully"
}

# 安装 protoc、protoc-gen-go
function iam::install::protobuf()
{
  # 检查 protoc、protoc-gen-go 是否安装
  command -v protoc &>/dev/null && command -v protoc-gen-go &>/dev/null && return 0

  # 1. 安装 protobuf
  rm -rf /tmp/protobuf

  cd /tmp/
  git clone --depth=1 https://github.com/protocolbuffers/protobuf
  cd protobuf
  ./autogen.sh
  ./configure
  make
  iam::common::sudo "make install"
  iam::log::info "install protoc tool successfully"

  # 2. 安装 protoc-gen-go
  echo $GO111MODULE
  go install github.com/golang/protobuf/protoc-gen-go@latest
  iam::log::info "install protoc-gen-go plugin successfully"
}

# ==================================================================