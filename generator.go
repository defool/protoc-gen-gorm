package main

import (
	"fmt"
	"strings"

	gorm "github.com/defool/protoc-gen-gorm/buf/generated/gorm/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type fieldInfo struct {
	name     string
	fileType string
	gormTag  string

	isRepeated bool
	isMessage  bool
}
type messageInfo struct {
	originType string
	ormType    string
	fields     []*fieldInfo
}

func generate(request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	plugin, err := protogen.Options{}.New(request)
	checkErr(err)

	for _, protoFile := range plugin.Files {
		if !isGormPB(protoFile) {
			continue
		}

		fileName := protoFile.GeneratedFilenamePrefix + ".pb.gorm.go"
		g := plugin.NewGeneratedFile(fileName, ".")

		var messages []messageInfo
		for _, message := range protoFile.Messages {
			if message.Desc.IsMapEntry() {
				continue
			}
			typeName := message.GoIdent.GoName
			messages = append(messages, messageInfo{
				originType: typeName,
				ormType:    typeName + "Orm",
				fields:     parseFields(message),
			})
		}

		g.P("package ", protoFile.GoPackageName)
		for _, m := range messages {
			generateMessageV2(g, m)
			generateConvertFunctionsV2(g, m)
		}

	}
	return plugin.Response(), nil
}

func parseFields(msg *protogen.Message) []*fieldInfo {
	ret := make([]*fieldInfo, 0)
	for _, field := range msg.Fields {
		fd := field.Desc
		options := fd.Options().(*descriptorpb.FieldOptions)
		fi := &fieldInfo{
			name:       field.GoName,
			isRepeated: field.Desc.Cardinality() == protoreflect.Repeated,
			isMessage:  field.Message != nil,
		}
		fieldType := fd.Kind().String()
		if v := proto.GetExtension(options, gorm.E_Gorm); v != nil {
			fi.gormTag = v.(string)
		}

		switch fieldType {
		case "float":
			fieldType = "float32"
		case "double":
			fieldType = "float64"
		}
		if field.Message != nil {
			parts := strings.Split(string(field.Desc.Message().FullName()), ".")
			fieldType = "*" + parts[len(parts)-1] + "Orm"
			if fi.isRepeated {
				fieldType = "[]" + fieldType
			}
		}
		fi.fileType = fieldType
		ret = append(ret, fi)
	}
	return ret
}

func generateConvertFunctionsV2(g *protogen.GeneratedFile, msg messageInfo) {
	///// To Orm
	g.P(`// ToOrm converts the pb message to orm message`)
	g.P(`func (m *`, msg.originType, `) ToOrm() *`, msg.ormType, ` {`)
	g.P(`if m == nil{return nil}`)
	g.P(`to := &`, msg.ormType, `{}`)
	for _, field := range msg.fields {
		generateFieldConversionV2(g, field, true)
	}
	g.P(`return to`)
	g.P(`}`)

	g.P()
	///// To Pb
	g.P(`// ToPb converts the orm message to pb message`)
	g.P(`func (m *`, msg.ormType, `) ToPb() *`, msg.originType, `{`)
	g.P(`if m == nil{return nil}`)
	g.P(`to := &`, msg.originType, `{}`)

	for _, field := range msg.fields {
		generateFieldConversionV2(g, field, false)
	}
	g.P(`return to`)
	g.P(`}`)
}

func generateMessageV2(g *protogen.GeneratedFile, m messageInfo) {
	g.P(`// `, m.ormType, ` is the ORM type for `, m.originType)
	g.P(`type `, m.ormType, ` struct {`)

	for _, f := range m.fields {
		tag := ""
		if f.gormTag != "" {
			tag = fmt.Sprintf("`gorm:\"%s\"`", f.gormTag)
		}
		g.P(f.name, ` `, f.fileType, tag)
	}
	g.P(`}`)
	g.P()
}

// Output code that will convert a field to/from orm.
func generateFieldConversionV2(g *protogen.GeneratedFile, field *fieldInfo, toORM bool) error {
	fieldName := field.name
	if field.isRepeated {
		g.P(`for _, v := range m.`, fieldName, ` {`)
		if toORM {
			g.P(`to.`, fieldName, ` = append(to.`, fieldName, `, v.ToOrm())`)
		} else {
			g.P(`to.`, fieldName, ` = append(to.`, fieldName, `, v.ToPb())`)
		}
		g.P(`}`)
	} else if field.isMessage {
		if toORM {
			g.P(`to.`, fieldName, ` = `, `m.`, fieldName, `.ToOrm()`)
		} else {
			g.P(`to.`, fieldName, ` = `,
				`m.`, fieldName, `.ToPb()`)
		}
	} else { // Singular raw ----------------------------------------------------
		g.P(`to.`, fieldName, ` = m.`, fieldName)
	}
	return nil
}

func isGormPB(protoFile *protogen.File) bool {
	for i := 0; i < protoFile.Desc.Imports().Len(); i++ {
		if protoFile.Desc.Imports().Get(i).Path() == "gorm/v1/gorm.proto" {
			return true
		}
	}
	return false
}
