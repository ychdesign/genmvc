package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"text/template"
)

type config struct {
	modelsPath string
	outputPath string
	tplPath    string
	fset       *token.FileSet
	tpl        *template.Template
	Components []Component
	goModule   string
}

func main() {
	c, err := NewConfig()
	if err != nil {
		fmt.Printf("new config err: %v\n", err)
		return
	}

	if c == nil {
		return
	}

	pkgs, err := c.parse()
	if err != nil {
		fmt.Printf("parse err: %v\n", err)
		return
	}

	err = c.validate()
	if err != nil {
		fmt.Printf("validate err: %v\n", err)
		return
	}

	err = c.process(pkgs)
	if err != nil {
		fmt.Printf("process err: %v\n", err)
		return
	}
}

// NewConfig 配置初始化
func NewConfig() (*config, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	args := os.Args[1:]
	var (
		// file flags
		all               = flag.Bool("all", false, "generate all. include bo、po、repostiroy、service and so on.")
		modelsPath        = flag.String("modelsPath", "pkg/models", "Path of the model's source code.")
		outputPath        = flag.String("outputPath", "generated", "Write all files to which directory")
		tplPath           = flag.String("fileTPLPath", home+"/.genmvc/templates", "Path of the templates to generate code.")
		POOption          = flag.Bool("po", false, "generate models to po entity.")
		BOOption          = flag.Bool("bo", false, "generate models to bo entity.")
		RepositoryOption  = flag.Bool("repository", false, "generate repository iface and implement.")
		RepositoryVersion = flag.String("repoVersion", "0.0.0", "repository template version.")
		ServiceOption     = flag.Bool("service", false, "generate service iface and implement.")
		ServiceVersion    = flag.String("svcVersion", "0.0.0", "service template version.")
	)

	flag.CommandLine.Parse(args)
	if flag.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return nil, nil
	}

	if *all {
		*POOption = true
		*BOOption = true
		*RepositoryOption = true
		*ServiceOption = true
	}

	cfg := &config{
		modelsPath: *modelsPath,
		outputPath: *outputPath,
		tplPath:    *tplPath,
		fset:       &token.FileSet{},
		Components: []Component{},
	}

	if *RepositoryOption {
		cfg.addComponent(NewRepositoryGenerator(cfg, "repositories", *RepositoryVersion))
	}

	if *ServiceOption {
		cfg.addComponent(NewServiceGenerator(cfg, "services", *ServiceVersion))
	}

	if *BOOption {
		cfg.addComponent(NewBOGenerator(cfg, "bo"))
	}

	if *POOption {
		cfg.addComponent(NewPOGenerator(cfg, "po"))
	}

	t := template.New("template")
	tpl, err := t.ParseGlob(*tplPath + "/*.tpl")
	if err != nil {
		return nil, err
	}
	cfg.tpl = tpl

	module, err := GetGoMod()
	if err != nil {
		return nil, err
	}
	cfg.goModule = module

	return cfg, nil
}

// validate 参数验证和处理
func (c *config) validate() error {
	os.Mkdir(c.outputPath, fs.ModePerm)

	return nil
}

// parse
func (c *config) parse() (map[string]*ast.Package, error) {
	c.fset = token.NewFileSet()
	pkgs, err := parser.ParseDir(c.fset, c.modelsPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

// process
func (c *config) process(pkgs map[string]*ast.Package) error {
	for _, element := range pkgs {
		for fileName, file := range element.Files {
			for _, component := range c.Components {
				if err := component.walk(fileName, file); err != nil {
					return err
				}
				if err := component.generate(fileName, file); err != nil {
					return err
				}
				if v, ok := component.(NodeReset); ok {
					v.reset(file)
				}
			}
		}
	}
	return nil
}

func (c *config) addComponent(cpt Component) {
	c.Components = append(c.Components, cpt)
}

// writeFile 输出到文件中
func writeFile(file string, body []byte) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.Write(body)

	if err != nil {
		return err
	}
	fmt.Printf("generate %v\n", file)
	return nil
}

// Component 代码生成器组件
type Component interface {
	walk(file string, node ast.Node) error
	generate(file string, node ast.Node) error
}

// NodeReset 复原源文件
type NodeReset interface {
	reset(node ast.Node)
}

// *****************************PO Generator*******************************

// poGenerator generate po model
type poGenerator struct {
	cfg   *config
	pkg   string
	Name  *ast.Ident
	Model string
}

func NewPOGenerator(cfg *config, pkg string) *poGenerator {
	return &poGenerator{cfg: cfg, pkg: pkg}
}

func (g *poGenerator) walk(file string, node ast.Node) error {
	fileNode, ok := node.(*ast.File)
	if !ok {
		return nil
	}
	g.Name = fileNode.Name
	tempIdent := *fileNode.Name
	fileNode.Name = &tempIdent
	tempIdent.Name = g.pkg

	ast.Inspect(fileNode, g.findModel)
	ast.Inspect(fileNode, g.rewrite)

	g.appendTableNameConst(fileNode)
	g.appendTableNamefunc(fileNode)
	return nil
}

func (g *poGenerator) appendTableNameConst(fileNode *ast.File) {
	tableNameConst := &ast.GenDecl{
		Doc: &ast.CommentGroup{},
		Tok: token.CONST,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{{
					Name: g.Model + "TableName",
				}},
				Values: []ast.Expr{&ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"\"",
				}},
				Comment: &ast.CommentGroup{},
			},
		},
	}
	fileNode.Decls = append(fileNode.Decls, tableNameConst)
}

func (g *poGenerator) appendTableNamefunc(fileNode *ast.File) {

	tableNamefunc := &ast.FuncDecl{
		Recv: &ast.FieldList{
			Opening: 0,
			List:    []*ast.Field{{Type: &ast.Ident{Name: g.Model}}},
			Closing: 0,
		},
		Name: &ast.Ident{
			Name: "TableName",
		},
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{Name: "string"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: g.Model + "TableName",
						},
					},
				}},
			Rbrace: 0,
		},
	}
	fileNode.Decls = append(fileNode.Decls, tableNamefunc)
}

func (g *poGenerator) generate(file string, node ast.Node) error {
	_, ok := node.(*ast.File)
	if !ok {
		return nil
	}

	buf := bytes.NewBuffer([]byte{})
	defer buf.Reset()

	err := format.Node(buf, g.cfg.fset, node)
	if err != nil {
		return err
	}

	os.Mkdir(path.Join(g.cfg.outputPath, g.pkg), fs.ModePerm)
	_, fileName := path.Split(file)

	return writeFile(path.Join(g.cfg.outputPath, g.pkg, fileName), buf.Bytes())
}

func (g *poGenerator) findModel(node ast.Node) bool {
	n, ok := node.(*ast.TypeSpec)
	if !ok {
		return true
	}
	if _, ok := n.Type.(*ast.StructType); ok {
		g.Model = n.Name.Name
	}
	return false
}

func (g *poGenerator) rewrite(node ast.Node) bool {

	x, ok := node.(*ast.StructType)
	if !ok {
		return true
	}

	for _, f := range x.Fields.List {

		fieldName := ""
		if len(f.Names) != 0 {
			for _, field := range f.Names {
				if isPublicName(field.Name) {
					fieldName = field.Name
					break
				}
			}
		}

		// nothing to process, continue with next line
		if fieldName == "" {
			continue
		}

		if f.Tag == nil {
			f.Tag = &ast.BasicLit{}
		}

		f.Tag.Value = quote(addGormTag(fieldName))
	}

	return false
}

func (g *poGenerator) reset(node ast.Node) {
	fileNode, ok := node.(*ast.File)
	if !ok {
		return
	}

	if g.Name != nil {
		fileNode.Name = g.Name
	}
}

// **************************** Repository Generator********************************

type repositoryGenerator struct {
	cfg     *config
	pkg     string
	Model   string
	version string
}

type repositoryModel struct {
	RepositoryIfaceName         string
	RepositoryIfaceInstanceName string
	RepositoryTableName         string
	Name                        string
	InstanceName                string
	InstanceSliceName           string
	ModulePath                  string
}

func NewRepositoryGenerator(cfg *config, pkg string, version string) Component {
	return &repositoryGenerator{
		cfg:     cfg,
		pkg:     pkg,
		version: version,
	}
}

func (g *repositoryGenerator) walk(file string, node ast.Node) error {
	fileNode, ok := node.(*ast.File)
	if !ok {
		return nil
	}
	ast.Inspect(fileNode, g.findModel)
	return nil
}

func (g *repositoryGenerator) generate(file string, node ast.Node) error {
	m := &repositoryModel{
		RepositoryIfaceName:         upperCaseFirst(g.Model) + "Repository",
		RepositoryIfaceInstanceName: lowerCaseFirst(g.Model) + "Repository",
		RepositoryTableName:         upperCaseFirst(g.Model) + "TableName",
		Name:                        upperCaseFirst(g.Model),
		InstanceName:                lowerCaseFirst(g.Model) + "PO",
		InstanceSliceName:           lowerCaseFirst(g.Model) + "POs",
		ModulePath:                  path.Join(g.cfg.goModule, g.cfg.outputPath),
	}

	bf := &bytes.Buffer{}
	err := g.cfg.tpl.ExecuteTemplate(bf, "repository_"+g.version+".tpl", m)
	if err != nil {
		return err
	}
	os.Mkdir(path.Join(g.cfg.outputPath, g.pkg), fs.ModePerm)
	_, fileName := path.Split(file)

	return writeFile(path.Join(g.cfg.outputPath, g.pkg, fileName), bf.Bytes())
}

func (g *repositoryGenerator) findModel(node ast.Node) bool {
	n, ok := node.(*ast.TypeSpec)
	if !ok {
		return true
	}
	if _, ok := n.Type.(*ast.StructType); ok {
		g.Model = n.Name.Name
	}

	return false
}

// ******************************BO Generator******************************

type boGenerator struct {
	cfg  *config
	pkg  string
	Name *ast.Ident
}

func NewBOGenerator(cfg *config, pkg string) Component {
	return &boGenerator{cfg: cfg, pkg: pkg}
}

func (g *boGenerator) walk(file string, node ast.Node) error {
	fileNode, ok := node.(*ast.File)
	if !ok {
		return nil
	}

	g.Name = fileNode.Name
	tempIdent := *fileNode.Name
	fileNode.Name = &tempIdent

	tempIdent.Name = g.pkg

	ast.Inspect(fileNode, g.rewrite)

	return nil
}

func (g *boGenerator) generate(file string, node ast.Node) error {
	_, ok := node.(*ast.File)
	if !ok {
		return nil
	}
	buf := bytes.NewBuffer([]byte{})
	defer buf.Reset()

	err := format.Node(buf, g.cfg.fset, node)
	if err != nil {
		return err
	}

	_, fileName := path.Split(file)
	os.Mkdir(path.Join(g.cfg.outputPath, g.pkg), fs.ModePerm)

	return writeFile(path.Join(g.cfg.outputPath, g.pkg, fileName), buf.Bytes())
}

func (g *boGenerator) reset(node ast.Node) {
	fileNode, ok := node.(*ast.File)
	if !ok {
		return
	}
	if g.Name != nil {
		fileNode.Name = g.Name
	}
}

func (g *boGenerator) rewrite(node ast.Node) bool {

	x, ok := node.(*ast.StructType)
	if !ok {
		return true
	}

	for _, f := range x.Fields.List {

		fieldName := ""
		if len(f.Names) != 0 {
			for _, field := range f.Names {
				if isPublicName(field.Name) {
					fieldName = field.Name
					break
				}
			}
		}

		// nothing to process, continue with next line
		if fieldName == "" {
			continue
		}

		if f.Tag == nil {
			f.Tag = &ast.BasicLit{}
		}

		f.Tag.Value = quote(addJsonTag(fieldName))
	}

	return false
}

// ******************************Service Generator******************************

type serviceGenerator struct {
	cfg     *config
	pkg     string
	Model   string
	version string
}

type serviceModel struct {
	ServiceIfaceName            string
	ServiceIfaceInstanceName    string
	RepositoryIfaceName         string
	RepositoryIfaceInstanceName string
	Name                        string
	InstanceName                string
	ModulePath                  string
}

func NewServiceGenerator(cfg *config, pkg string, version string) Component {
	return &serviceGenerator{
		cfg:     cfg,
		pkg:     pkg,
		version: version,
	}
}

func (g *serviceGenerator) walk(file string, node ast.Node) error {
	fileNode, ok := node.(*ast.File)
	if !ok {
		return nil
	}
	ast.Inspect(fileNode, g.findModel)
	return nil
}

func (g *serviceGenerator) generate(file string, node ast.Node) error {
	m := &serviceModel{
		ServiceIfaceName:            upperCaseFirst(g.Model) + "Service",
		ServiceIfaceInstanceName:    lowerCaseFirst(g.Model) + "Service",
		RepositoryIfaceName:         upperCaseFirst(g.Model) + "Repository",
		RepositoryIfaceInstanceName: lowerCaseFirst(g.Model) + "Repository",
		Name:                        upperCaseFirst(g.Model),
		InstanceName:                lowerCaseFirst(g.Model) + "PO",
		ModulePath:                  path.Join(g.cfg.goModule, g.cfg.outputPath),
	}

	bf := &bytes.Buffer{}
	err := g.cfg.tpl.ExecuteTemplate(bf, "service_"+g.version+".tpl", m)
	if err != nil {
		return err
	}

	os.Mkdir(path.Join(g.cfg.outputPath, g.pkg), fs.ModePerm)
	_, fileName := path.Split(file)

	return writeFile(path.Join(g.cfg.outputPath, g.pkg, fileName), bf.Bytes())
}

func (g *serviceGenerator) findModel(node ast.Node) bool {
	n, ok := node.(*ast.TypeSpec)
	if !ok {
		return true
	}
	if _, ok := n.Type.(*ast.StructType); ok {
		g.Model = n.Name.Name
	}

	return false
}
