GOPATH:=~/go
PYTHONPATH=.
PYTHON=python3
PY2CHROMA=$(GOPATH)/src/github.com/alecthomas/chroma/_tools/pygments2chroma.py
INPUT_FILE=sMark.sMarkLexer
INPUT_NO_PYTHON=highlight/sMark.chroma.go
OUTPUT_FILE=$(GOPATH)/src/github.com/alecthomas/chroma/lexers/sMark.go
OUTPUT_FILE_LOCAL=$(INPUT_NO_PYTHON)


#Rebuild chroma adding sMark lexer
chroma:
	PYTHONPATH=$(PYTHONPATH) $(PYTHON) $(PY2CHROMA) $(INPUT_FILE) > $(OUTPUT_FILE)
	go install github.com/alecthomas/chroma github.com/alecthomas/chroma/lexers github.com/alecthomas/chroma/cmd/chroma

#Compile the sMark lexer but not install it
chroma-compile:
	PYTHONPATH=$(PYTHONPATH) $(PYTHON) $(PY2CHROMA) $(INPUT_FILE) > $(OUTPUT_FILE_LOCAL)

#Rebuild chroma with the precompiled version of sMark lexer
chroma-no-python:
	cp $(INPUT_NO_PYTHON) $(OUTPUT_FILE)
	go install github.com/alecthomas/chroma github.com/alecthomas/chroma/lexers github.com/alecthomas/chroma/cmd/chroma