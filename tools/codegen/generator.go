package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"golang.org/x/tools/go/packages"
	"html/template"
	"log"
	"strings"
)

// Generator 保存分析的状态。主要用于缓冲 format.Source 的输出。
type Generator struct {
	buf bytes.Buffer // Accumulated output(累计输出).
	pkg *Package     // Package we are scanning（我们正在扫描的包）.

	trimPrefix string
}

// Printf 与 fmt.Printf 类似，但将字符串添加到 g.buf。
func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// parsePackage 分析由 patterns 和 tags 构造的单个包。
// 如果出现错误，则 parsePackage 退出。
func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		// nolint: staticcheck
		Mode: packages.LoadSyntax,
		Tests: false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

// addPackage 将类型检查包及其语法文件添加到生成器。
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file:       file,
			pkg:        g.pkg,
			trimPrefix: g.trimPrefix,
		}
	}
}

// generate produces the register calls for the named type.
func (g *Generator) generate(typeName string) {
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}
	// Generate code that will fail if the constants change value.
	g.Printf("\t// init register error codes defines in this source code to `github.com/marmotedu/errors`\n")
	g.Printf("func init() {\n")
	for _, v := range values {
		code, description := v.ParseComment()
		g.Printf("\tregister(%s, %s, \"%s\")\n", v.originalName, code, description)
	}
	g.Printf("}\n")
}

// generateDocs produces error code markdown document for the named type.
func (g *Generator) generateDocs(typeName string) {
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}

	tmpl, _ := template.New("doc").Parse(errCodeDocPrefix)
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, "`")

	// Generate code that will fail if the constants change value.
	g.Printf(buf.String())
	for _, v := range values {
		code, description := v.ParseComment()
		// g.Printf("\tregister(%s, %s, \"%s\")\n", v.originalName, code, description)
		g.Printf("| %s | %d | %s | %s |\n", v.originalName, v.value, code, description)
	}
	g.Printf("\n")
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")

		return g.buf.Bytes()
	}

	return src
}
