package main

import (
	"fmt"
	"go/ast"
	"strings"

	gorm "github.com/defool/protoc-gen-gorm/buf/generated/gorm/v1"
	"github.com/fatih/structtag"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type myVisitor struct {
	pgs.Visitor
	pgs.DebuggerCommon
	pgsgo.Context
	pctx    pgsgo.Context
	gormTag map[string]string
}

func (v *myVisitor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	var gTag string
	ok, err := f.Extension(gorm.E_Gorm, &gTag)
	if err != nil {
		return nil, err
	}
	if !ok {
		return v, nil
	}
	msgName := v.pctx.Name(f.Message()).String()
	fieldName := v.pctx.Name(f).String()
	key := fmt.Sprintf("%s/%s", msgName, fieldName)
	v.gormTag[key] = gTag
	return v, nil
}

func getTags(f pgs.File, m *mod) (map[string]string, error) {
	v := &myVisitor{DebuggerCommon: m, Context: m.Context, pctx: m.Context, gormTag: make(map[string]string)}
	v.Visitor = pgs.PassThroughVisitor(v)
	err := pgs.Walk(v, f)
	return v.gormTag, err
}

func retagOnAst(n ast.Node, tags map[string]string) error {
	r := tagVisitor{tags: tags}
	ast.Walk(r, n)
	return r.err
}

type tagVisitor struct {
	tags     map[string]string
	err      error
	typeName string
}

func (v tagVisitor) Visit(n ast.Node) ast.Visitor {
	if tp, ok := n.(*ast.TypeSpec); ok {
		v.typeName = tp.Name.String()
		logger.Println("type is", v.typeName)
		return v
	}

	f, ok := n.(*ast.Field)
	if !ok {
		return v
	}
	if len(f.Names) == 0 {
		return nil
	}

	fieldName := f.Names[0].String()
	// Only for exportable field
	if len(fieldName) == 0 || (fieldName[0] > 'Z' || fieldName[0] < 'A') {
		return nil
	}

	var allTagVal string
	if f.Tag != nil {
		allTagVal = strings.Trim(f.Tag.Value, "`")
	}
	oldTags, err := structtag.Parse(allTagVal)
	if err != nil {
		logger.Println("parse old tag error", err)
		return nil
	}

	key := fmt.Sprintf("%s/%s", v.typeName, fieldName)
	pbTag := v.tags[key]

	var hasColumn bool
	if pbTag != "" {
		for _, v := range strings.Split(pbTag, ";") {
			if strings.HasPrefix(v, "column:") {
				hasColumn = true
			}
		}
	}
	if !hasColumn && pbTag != "-" {
		if pbTag != "" && !strings.HasSuffix(pbTag, ";") {
			pbTag += ";"
		}
		columnName := camelCaseToUnderscore(fieldName)
		columnName = applyReplaceKeyword(v.typeName, columnName)
		// only for none-pointer type
		if !strings.HasPrefix(fmt.Sprint(f.Type), "&") {
			pbTag += "column:" + columnName
		}
	}

	// replace
	if pbTag != "" {
		oldTags.Set(&structtag.Tag{Key: "gorm", Name: pbTag})
	}

	f.Tag.Value = "`" + oldTags.String() + "`"
	return nil
}

func applyReplaceKeyword(typeName, s string) string {
	if opt.replaceKeyword && keywordMap[s] {
		return camelCaseToUnderscore(typeName) + "_" + s
	}
	return s
}
