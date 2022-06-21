package main

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
	"path/filepath"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type mod struct {
	*pgs.ModuleBase
	pgsgo.Context
}

func newMod() pgs.Module {
	return &mod{ModuleBase: &pgs.ModuleBase{}}
}

func (m *mod) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

func (mod) Name() string {
	return "gorm"
}

func (m mod) Execute(targets map[string]pgs.File, gpkgs map[string]pgs.Package) []pgs.Artifact {
	applyOptions(m.Parameters().String())
	initLogger(enableLogger)

	logger.Printf("opts is %+v | %s\n", opt, m.Parameters().String())

	for _, f := range targets {
		if !isGormPB(f.Imports()) {
			continue
		}
		m.retag(f)
		m.genGorm(f)

	}

	return m.Artifacts()
}

func (m mod) genGorm(f pgs.File) {
	pkgName := *f.Descriptor().Package
	if idx := strings.LastIndexByte(pkgName, '.'); idx >= 0 {
		pkgName = pkgName[idx+1:]
	}
	sourceFile := *f.Descriptor().Name
	name := strings.Title(filepath.Base(sourceFile))
	if idx := strings.IndexByte(name, '.'); idx >= 0 {
		name = name[:idx]
	}
	info := FileFieldInfo{
		Name:      name,
		Source:    sourceFile,
		Package:   pkgName,
		GenFields: opt.generateField,
	}
	for _, msg := range f.AllMessages() {
		typeName := underscoreToCamelCase(*msg.Descriptor().Name)
		var fields []FieldInfo
		for _, fe := range msg.Fields() {
			columnName := camelCaseToUnderscore(*fe.Descriptor().Name)
			goName := underscoreToCamelCase(columnName)
			columnName = applyReplaceKeyword(typeName, columnName)
			fields = append(fields, FieldInfo{Name: goName, Field: columnName})
		}

		info.Messages = append(info.Messages, MessageFieldInfo{
			Name:   typeName,
			Fields: fields,
		})
	}

	buf := bytes.NewBuffer(nil)
	tp := template.Must(template.New("").Parse(fieldTemplate))
	err := tp.Execute(buf, info)
	m.CheckErr(err)

	filename := m.Context.OutputPath(f).SetExt(".gorm.go").String()
	if opt.outdir != "" {
		filename = filepath.Join(opt.outdir, filename)
	}

	m.AddGeneratorFile(filename, buf.String())
}

func (m mod) retag(f pgs.File) {
	tags, err := getTags(f, &m)
	logger.Println("imports", f.Imports())
	logger.Println("tags is", tags, err)
	m.CheckErr(err)
	filename := m.Context.OutputPath(f).SetExt(".go").String()
	if opt.outdir != "" {
		filename = filepath.Join(opt.outdir, filename)
	}
	fs := token.NewFileSet()
	fn, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	m.CheckErr(err)
	err = retagOnAst(fn, tags)
	m.CheckErr(err)

	var buf strings.Builder
	m.CheckErr(printer.Fprint(&buf, fs, fn))

	m.OverwriteGeneratorFile(filename, buf.String())
}

func isGormPB(files []pgs.File) bool {
	for i := 0; i < len(files); i++ {
		if *files[i].Descriptor().Name == "gorm/v1/gorm.proto" {
			return true
		}
	}
	return false
}
