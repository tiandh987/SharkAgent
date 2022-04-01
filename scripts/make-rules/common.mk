
# MAKEFILE_LIST 特殊变量：
# 	make 所需要处理的 Makefile 文件列表
#	当前 makefile 的文件名总是位于列表的最后，文件名之间以“空格”进行分割
#
# $(lastword <text>) 预定义函数：
#	取字符串 <text> 中的最后一个单词，并返回字符串 <text> 的最后一个单词
#
# $(dir <names...>) 预定义函数
# 	从文件名序列中取出目录部分。
#	目录部分是指最后一个反斜杠（/）之前的部分。
#	返回文件名序列的非目录部分。
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

# $(origin <variable>) 预定义函数
# 	告诉变量的 "出生情况", 有如下返回值:
#		undefined: <variable> 从来没有定义过
#
# $(shell cat foo) 预定义函数
# 	执行操作系统命令, 并返回操作结果
#
# $(abspath <text>) 预定义函数
# 	将 <text> 中的各路径转换成绝对路径, 并将转换后的结果返回
#
# pwd -P
# 	打印物理目录，没有任何符号链接
ifeq ($(origin ROOT_DIR), undefined)
ROOT_DIR := $(abspath $(shell cd ${COMMON_SELF_DIR}/../.. && pwd -P))
endif