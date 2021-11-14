default: help

## help : This help
help: Makefile
		@printf "\n Dreamstation\n\n"
		@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
		@printf ""

## install : [INSTALL] Crée le virtualenv et installe les paquets python
install: venv install_pip

## install_pip : [INSTALL] Installe les paquets python au sein du venv
install_pip:
	./venv/bin/pip install -r ./requirements.txt

## venv : [INSTALL] Crée un virtualenv vierge
venv:
	virtualenv venv --python=`which python3.7`

## req : [DEV] Freeze pyhton requirements
req:
	./venv/bin/pip freeze | grep -v "pkg-resources" > requirements.txt
