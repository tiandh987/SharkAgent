#!/usr/bin/env bash

# 本脚本功能：根据 scripts/environment.sh 配置，生成 IAM 组件 YAML 配置文件。
# 示例：genconfig.sh scripts/environment.sh configs/iam-apiserver.yaml

env_file="$1"
template_file="$2"

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${IAM_ROOT}/scripts/lib/init.sh"

if [ $# -ne 2 ];then
    iam::log::error "Usage: genconfig.sh scripts/environment.sh configs/iam-apiserver.yaml"
    exit 1
fi

source "${env_file}"

# declare 设置变量值和属性。
#   -A	to make NAMEs associative arrays (if supported)
declare -A envs

# set 设置或取消设置 shell 选项和位置参数的值。
# 更改 shell 属性和位置参数的值，或显示 shell 变量的名称和值。
#
# 使用 + 而不是 - 会导致这些标志被关闭。
# 这些标志也可以在调用 shell 时使用。
# 当前的标志集可以在 $- 中找到。
# 剩余的 n 个 ARGs 是位置参数，并按顺序分配给 $1、$2、.. $n。
# 如果没有给出 ARG，则打印所有 shell 变量。
#
# -u 替换时将未设置的变量视为错误。
set +u

# sed -n '/oo/p'
#   p: 数据的搜寻并显示
#
# sed 之 \1-9 的作用：
#   https://blog.csdn.net/jasonliujintao/article/details/53509732
#   https://cloud.tencent.com/developer/ask/113794
#
# ^[^#].*${\(.*\)}.*
#   匹配非 # 字符开始，包含 ${xxx} 的字符
#   \1 指的就是 xxx 被 () 包含的部分
#
# eval
#   eval: eval [arg ...]
#   将参数作为 shell 命令执行。
#   将 ARGs 组合成一个字符串，将结果用作 shell 的输入，然后执行结果命令。
#
#   https://chegva.com/3818.html#:~:text=eval,%E7%94%A8%E6%9D%A5%E5%9C%A8%E6%89%A7%E8%A1%8C%E5%91%BD%E4%BB%A4%E6%97%B6%E4%BD%9C%E4%BA%8C%E6%AC%A1%E8%A7%A3%E6%9E%90%EF%BC%9A%E4%B8%BB%E8%A6%81%E6%98%AF%E6%AF%8F%E6%AC%A1%E6%89%A7%E8%A1%8C%E4%B8%80%E4%B8%AAshell%E5%91%BD%E4%BB%A4%E5%AE%83%E4%BC%9A%E5%85%88%E6%A3%80%E5%AF%9F%E4%B8%80%E6%AC%A1%EF%BC%8C%E7%9C%8B%E5%88%B0%E6%9C%89%24%E6%A0%87%E5%BF%97%E5%B0%B1%E4%BC%9A%E6%8A%8A%E5%80%BC%E6%9B%BF%E6%8D%A2%E4%B8%80%E6%AC%A1%EF%BC%8C%E7%84%B6%E5%90%8E%E5%86%8D%E6%89%A7%E8%A1%8C%E4%B8%80%E9%81%8D%E3%80%82
#   eval用来在执行命令时作二次解析：
#     主要是每次执行一个shell命令它会先检察一次，看到有$标志就会把值替换一次，然后再执行一遍。
#
# -z
#   [ -z STRING ] 如果STRING的长度为零则为真 ，即判断是否为空，空即是真；
#   https://www.cnblogs.com/liudianer/p/12071476.html
#
# 检查配置文件模板中的定义变量 xxx（${xxx}）是否被设置
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file})
do
    if [ -z "$(eval echo \$${env})" ];then
        iam::log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ "${missing}" ];then
    iam::log::error 'You may run `source scripts/environment.sh` to set these environment'
fi

# 通过在 environment.sh 中定义环境变量 foo=bar
# 配置文件模板中定义 ${foo}
# eval 会在执行命令时二次解析，将 ${xxx} 解析为 bar，从而生成配置文件
eval "cat << EOF
$(cat ${template_file})
EOF"