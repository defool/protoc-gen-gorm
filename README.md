# protoc-gen-gorm

`protoc-gen-gorm` is a protoc plugin for injecting `[gorm](https//gorm.io)` tag to protobuf message。

[中文介绍](./README.zh.md)

## How it works

Inject tag by `github.com/lyft/protoc-gen-star(github.com/lyft/protoc-gen-star)` after plugin `[protoc-gen-go](https://golang/protobuf/protoc-gen-go)`.

## Example
 
 ./example/foo/v1/db.proto: 
```
syntax = "proto3";
package foo.v1;
option go_package="foo/v1";

// reference of proto file in `./buf`
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

generate stub file in two step by protoc：
```
protoc  -I . -I ./buf  --go_out="./example/generated"  ./example/foo/v1/db.proto
protoc  -I . -I ./buf  --gorm_out="outdir=./example/generated:."  ./example/foo/v1/db.proto
# Note: The out argment is replace by outdir option
```

OR `buf` use case:

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

generate stub file in two step：
```
buf generate
buf generate --template buf.gen.gorm.yaml
```

## Other features

- When use the option `replace_keyword=true`，the column name will be replace by table+column if it's keyword in MySQL.
- Generate the database column name as variable value to avoid using column name in code directly.