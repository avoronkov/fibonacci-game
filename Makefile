.PHONY: all clean test html-game

fibonacci-game: test
	go build

test:
	cd common && \
	go test

html/jsgame.js:
	cd html && \
	gopherjs build -m jsgame.go

html/fibonacci-game.html: html/jsgame.js html/fibonacci-game.html.tpl
	cd html && \
	perl -nle 'if (/@@#include\s+"([^"]+)"/) { print `cat "$$1"` } else { print }' fibonacci-game.html.tpl > fibonacci-game.html

html-game: html/fibonacci-game.html

all: fibonacci-game html-game

clean:
	rm -f fibonacci-game html/fibonacci-game.html html/jsgame.js html/jsgame.js.map
