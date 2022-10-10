# genmvc

golang mvc code-generator.

## Install

```bash
git clone git@github.com:ychdesign/genmvc.git
cd genmvc
make install
```

## Usage

exec ``genmvc`` in terminal we can get:

```
Usage of genmvc:
  -all
    	generate all. include bo、po、repostiroy、service and so on.
  -bo
    	generate models to bo entity.
  -fileTPLPath string
    	Path of the templates to generate code. (default "~/.genmvc/templates")
  -modelsPath string
    	Path of the models source code. (default "pkg/models")
  -outputPath string
    	Write all file to which directory (default "generated")
  -po
    	generate models to po entity.
  -repoVersion string
    	repository template version. (default "0.0.0")
  -repository
    	generate repository iface and implement.
  -service
    	generate service iface and implement.
  -svcVersion string
    	service template version. (default "0.0.0")
```

### Examples

```bash
$ genmvc -all -modelsPath examples/models -outputPath examples/generated
$ go mod tidy
```

Output:

```sh
generate examples/generated/repositories/server.go
generate examples/generated/services/server.go
generate examples/generated/bo/server.go
generate examples/generated/po/server.go
```

```bash
$tree examples

examples/
├── generated
│   ├── bo
│   │   └── server.go
│   ├── po
│   │   └── server.go
│   ├── repositories
│   │   └── server.go
│   └── services
│       └── server.go
└── models
    └── server.go

6 directories, 5 files
```
