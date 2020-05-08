FROM golang:1.14-buster as build-env

WORKDIR /app
COPY . /app

RUN make all

FROM alpeware/chrome-headless-trunk

WORKDIR /app

COPY --from=build-env /app/bin/*                  /app/bin/
COPY --from=build-env /app/Makefile               /app/
COPY --from=build-env /app/configuration/defaults /app/configuration/defaults

CMD ["/app/bin/roastmap"]