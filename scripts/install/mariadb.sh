#!/usr/bin/env bash

# 项目代码根目录
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

# 空变量和没有初始化的变量可能会对shell脚本测试产生灾难性的影响，
# 因此在不确定变量的内容的时候，在测试号前使用 -n 或者 -z 测试一下
#
# [ -z STRING ] 如果STRING的长度为零则返回为真，即空是真
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# 安装
function iam::mariadb::install()
{
  # 1. 配置 MariaDB 10.5 Yum 源
  echo ${LINUX_PASSWORD} | sudo -S bash -c "cat << 'EOF' > /etc/yum.repos.d/mariadb-10.5.repo
# MariaDB 10.5 CentOS repository list - created 2020-10-23 01:54 UTC
# http://downloads.mariadb.org/mariadb/repositories/
[mariadb]
name = MariaDB
baseurl = https://mirrors.aliyun.com/mariadb/yum/10.5/centos8-amd64/
module_hotfixes=1
gpgkey=https://yum.mariadb.org/RPM-GPG-KEY-MariaDB
gpgcheck=0
EOF"

  # 2. 安装MariaDB和MariaDB客户端
  iam::common::sudo "yum -y install MariaDB-server MariaDB-client"

  # 3. 启动 MariaDB，并设置开机启动
  iam::common::sudo "systemctl enable mariadb"
  iam::common::sudo "systemctl start mariadb"

  # 4. 设置root初始密码
  iam::common::sudo "mysqladmin -u${MARIADB_ADMIN_USERNAME} password ${MARIADB_ADMIN_PASSWORD}"

  iam::mariadb::status || return 1
  iam::mariadb::info
  iam::log::info "install MariaDB successfully"
}

# 卸载
function iam::mariadb::uninstall()
{
  set +o errexit
  iam::common::sudo "systemctl stop mariadb"
  iam::common::sudo "systemctl disable mariadb"
  iam::common::sudo "yum -y remove MariaDB-server MariaDB-client"
  iam::common::sudo "rm -rf /var/lib/mysql"
  iam::common::sudo "rm -f /etc/yum.repos.d/mariadb-10.5.repo"
  set -o errexit
  iam::log::info "uninstall MariaDB successfully"
}

# 安装后打印必要的信息
function iam::mariadb::info() {
cat << EOF
MariaDB Login: mysql -h127.0.0.1 -u${MARIADB_ADMIN_USERNAME} -p'${MARIADB_ADMIN_PASSWORD}'
EOF
}

# 状态检查
function iam::mariadb::status()
{
  # 查看mariadb运行状态，如果输出中包含active (running)字样说明mariadb成功启动。
  systemctl status mariadb |grep -q 'active' || {
    iam::log::error "mariadb failed to start, maybe not installed properly"
    return 1
  }

  mysql -u${MARIADB_ADMIN_USERNAME} -p${MARIADB_ADMIN_PASSWORD} -e quit &>/dev/null || {
    iam::log::error "can not login with root, mariadb maybe not initialized properly"
    return 1
  }
}

# Shell编程中，我们可以使用双中括号运算符 [[]] 和 =~ 来判断字符串是否匹配给定的正则表达式
# $* 传递给脚本或函数的所有参数。
if [[ "$*" =~ iam::mariadb:: ]];then
  eval $*
fi
