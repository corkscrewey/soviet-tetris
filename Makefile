SHELL=/bin/sh

run: simh mame
.PHONY: run

simh:
	docker compose up -d
.PHONY: simh

mame:
	./mame/mame ie15 -rompath files/rom -window -rs232 null_modem -bitb socket.localhost:2323
.PHONY: mame

build:
	docker compose build
.PHONY: build
