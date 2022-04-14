#!/usr/bin/env bash

# 项目代码根目录
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..

# 空变量和没有初始化的变量可能会对shell脚本测试产生灾难性的影响，
# 因此在不确定变量的内容的时候，在测试号前使用 -n 或者 -z 测试一下
#
# [ -z STRING ] 如果STRING的长度为零则返回为真，即空是真
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# 安装
function iam::redis::install()
{
  # 1. 安装 Redis
  iam::common::sudo "yum -y install redis"

  # 2. 配置 Redis
  # 2.1 修改`/etc/redis.conf`文件，将 daemonize 由 no 改成 yes，表示允许 redis 在后台启动
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^daemonize/{s/no/yes/}' /etc/redis.conf

  # 2.2 在`bind 127.0.0.1`前面添加 `#` 将其注释掉，默认情况下只允许本地连接，注释掉后外网可以连接Redis
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^# bind 127.0.0.1/{s/# //}' /etc/redis.conf

  # 2.3 修改 requirepass 配置，设置Redis密码
  echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^# requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' /etc/redis.conf

  # 2.4 因为我们上面配置了密码登录，需要将 protected-mode 设置为 no，关闭保护模式
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^protected-mode/{s/yes/no/}' /etc/redis.conf

  # 3. 为了能够远程连上Redis，需要执行以下命令关闭防火墙，并禁止防火墙开机启动（如果不需要远程连接，可忽略此步骤）
  iam::common::sudo "systemctl stop firewalld.service"
  iam::common::sudo "systemctl disable firewalld.service"

  # 4. 启动 Redis
  iam::common::sudo "redis-server /etc/redis.conf"

  iam::redis::status || return 1
  iam::redis::info
  iam::log::info "install Redis successfully"
}

# 卸载
function iam::redis::uninstall()
{
  set +o errexit
  iam::common::sudo "killall redis-server"
  iam::common::sudo "yum -y remove redis"
  iam::common::sudo "rm -rf /var/lib/redis"
  set -o errexit
  iam::log::info "uninstall Redis successfully"
}

# 安装后打印必要的信息
function iam::redis::info() {
cat << EOF
Redis Login: redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a '${REDIS_PASSWORD}'
EOF
}

# 状态检查
function iam::redis::status()
{
  if [[ -z "`pgrep redis-server`" ]];then
    iam::log::error_exit "Redis not running, maybe not installed properly"
    return 1
  fi


  redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a "${REDIS_PASSWORD}" --hotkeys || {
    iam::log::error "can not login with ${REDIS_USERNAME}, redis maybe not initialized properly"
    return 1
  }
}

# Shell编程中，我们可以使用双中括号运算符 [[]] 和 =~ 来判断字符串是否匹配给定的正则表达式
# $* 传递给脚本或函数的所有参数。
if [[ "$*" =~ iam::redis:: ]];then
  eval $*
fi


