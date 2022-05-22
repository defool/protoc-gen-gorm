package main

import (
	"html/template"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	val := FileFieldInfo{
		Source:  "foo/v1/db.proto",
		Package: "v1",
		Name:    "Db",
		Messages: []MessageFieldInfo{
			{
				Name: "User",
				Fields: []FieldInfo{
					{
						Name:  "Email",
						Field: "email",
					},
					{
						Name:  "Name",
						Field: "user_name",
					},
				},
			},
			{
				Name: "Group",
				Fields: []FieldInfo{
					{
						Name:  "Name",
						Field: "group_name",
					},
				},
			},
		},
	}
	tp := template.Must(template.New("").Parse(fieldTemplate))
	err := tp.Execute(os.Stdout, val)
	checkErr(err)
}
