FROM alpine:3.11.3

COPY kubeyaml /kubeyaml

ENTRYPOINT ["/kubeyaml"]