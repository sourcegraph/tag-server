.PHONY: install clean

install: .bin/srclib-ctags

clean:
	rm .bin/srclib-ctags

.bin/srclib-ctags:
	go build -o .bin/srclib-ctags
