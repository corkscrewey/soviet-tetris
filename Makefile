SHELL=/bin/sh

ifndef MAME
MAME=./mame/mame
endif

run: simh mame
.PHONY: run

simh:
	docker compose up -d
.PHONY: simh

mame:
	$(MAME) ie15 -rompath files/rom -window -rs232 null_modem -bitb socket.localhost:2323
.PHONY: mame

build:
	docker compose build
.PHONY: build
