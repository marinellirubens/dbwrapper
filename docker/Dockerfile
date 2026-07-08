FROM golang:1.23.4

RUN mkdir -p /tmp/dbwrapper

ARG VERSION_STR

COPY ./bin/${VERSION_STR}/dbwrapper /bin/dbwrapper
COPY ./internal/config/examples/config.example.json /opt/config.json

CMD ["/bin/dbwrapper", "-f", "/opt/config.json"]

