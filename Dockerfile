FROM golang:1.17-buster as build-env

WORKDIR /app
COPY . /app

RUN make all

FROM alpeware/chrome-headless-trunk

WORKDIR /app

RUN apt-get update -qqy \
 && apt-get -qqy install curl \
 && rm -rf /var/lib/apt/lists/* /var/cache/apt/*

COPY --from=build-env /app/bin/*                  /app/bin/
COPY --from=build-env /app/Makefile               /app/
COPY --from=build-env /app/configuration/defaults /app/configuration/defaults

CMD ["/app/bin/prerender run"]