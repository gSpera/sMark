INPUT_FILE=sMark.sm
OUTPUT_FILE=../index.html
SCSS_INPUT=scss/index.scss
SCSS_OUTPUT=../index.css

compile:
	scss $(SCSS_INPUT):$(SCSS_OUTPUT)
	sMark -i $(INPUT_FILE) -o $(OUTPUT_FILE)

clean:
	rm ../index.*