
# Makefile 允许对目标进行类似正则运算的匹配，主要用到的通配符是%。
# 通过使用通配符，可以使不同的目标使用相同的规则，从而使 Makefile 扩展性更强，也更简洁。
#
# $*
#	这个变量表示目标模式中 % 及其之前的部分。
#	如果目标是 dir/a.foo.b ，并且目标的模式是 a.%.b ,
#	那么, $* 的值就会是 dir/a.foo（？？？） 实际得到的是 dir/foo
#	示例：a.%.b            dir/a.foo.b            -->  dir/foo
#		 a.%.b            dir1/dir2/a.foo.b      -->  dir1/dir2/foo
#		 a.%.b            dir1/dir2/aa.foo.b     -->  报错
#		 a.%.b            dir1/dir2/a.foo.bb     -->  报错
#		 tools.install.%  tools.install.codegen  -->  codegen
#
# $(MAKE) 特殊变量
#	make
.PHONE: tools.install.%
tools.install.%:
	echo "==========> Installing $*"
	$(MAKE) install.$*

# $(GO) golang.mk 中定义
#	go
#
# ROOT_DIR common.mk 中定义
#
# 代码生成工具安装；
#	使用：在 SharkAgent 根目录下，执行 make tools.install.codegen
.PHONE: install.codegen
install.codegen:
	$(GO) install ${ROOT_DIR}/tools/codegen/codegen.go