FROM alpine:3.17 as build
LABEL maintainer="github.com/corkscrewey"
RUN apk add --no-cache gcc musl-dev curl make

WORKDIR /build
RUN curl -LO https://github.com/simh/simh/archive/refs/tags/v3.9-0.tar.gz \
	&& tar xf v3.9-0.tar.gz
RUN cd simh-3.9-0 && make BIN/pdp11 

# get the disk images from github "files" release
RUN curl -LO https://github.com/corkscrewey/soviet-tetris/releases/download/files/disks.zip \
	&& unzip disks.zip

FROM alpine:3.17

COPY --from=build /build/simh-3.9-0/BIN/pdp11 /usr/local/bin/

WORKDIR /sim
COPY --from=build /build/*.dsk /sim/
COPY files/config/pdp.ini .

EXPOSE 2323

CMD ["/usr/local/bin/pdp11", "-i", "pdp.ini"]
