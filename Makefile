help:
	@echo "Makefile rules:"
	@echo
	@grep -E '^### [-A-Za-z0-9_]+:' Makefile | sed 's/###/   /'

### install: install cmd genmvc to $GOBIN
install:
	@echo "init genmvc"
	@mkdir -p ~/.genmvc
	@cp -r ./pkg/templates ~/.genmvc
	go install .

### update: update templates
update: 
	@echo "update templates"
	@cp -r ./pkg/templates ~/.genmvc
	go install .
