PROG_NAME = ep

.PHONY: all clean $(PROG_NAME)

$(PROG_NAME): static
	go build .

emojis:
	mkdir -p static
	curl -o static/emojis.json https://raw.githubusercontent.com/muan/emojilib/master/emojis.json

static: emojis
	go generate .