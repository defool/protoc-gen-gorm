# protoc-gen-gorm

`protoc-gen-gorm`是用来在Protobuf Message结构体中注入 [gorm](https//gorm.io) 标签(Tag)的 `protoc` 插件。

[English](./README.md)

## 原理

参考了[protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag)的实现，在 `protoc-gen-go`插件生成代码后，再使用`github.com/lyft/protoc-gen-star`注入`gorm`的Tag。

## 安装

```
go install github.com/defool/protoc-gen-gorm
```

## 示例
 
 ./example/foo/v1/db.proto: 
```
syntax = "proto3";
package foo.v1;
option go_package="foo/v1";

// 引用buf目录中的proto文件
import "gorm/v1/gorm.proto";

message User {
    uint64 id = 1;
    string name = 2 [(gorm)="size:32;column:uname;"];
    string user_email = 3 [(gorm)="size:32;"];
    uint64 company_id = 4;
    Company company = 5;
    repeated Group groups = 6 [(gorm)="many2many:user_languages;"]; 
}

message Company {
    uint64 id = 1;
    string name = 2;
}

message Group {
    uint64 id = 1;
    string name = 2;
}
```

使用protc生成桩代码:
```
protoc  -I . -I ./buf  --go_out="./example/generated"  ./example/foo/v1/db.proto
protoc  -I . -I ./buf  --gorm_out="outdir=./example/generated:."  ./example/foo/v1/db.proto
# 注意gorm插件的输出的目录通过outdir参数传入
```

或使用`buf`来生成

buf.gen.yaml:
```
version: v1
plugins:
  - name: go
    out: generated
    opt: paths=source_relative
```

buf.gen.gorm.yaml:
```
version: v1
plugins:  
  - name: gorm
    out: .
    opt:
    - paths=source_relative
    - outdir=./generated
    - replace_keyword=true
```

需要执行二步来生成桩代码：
```
buf generate
buf generate --template buf.gen.gorm.yaml
```

## 其他功能

- 插件选项`replace_keyword=true`时，如果gorm对应的column名是MySQL的关键字，则使用表名_原字段名替换原字段名
- 生成gorm的字段名常量，避免在代码中直接使用数据库中的字段名