FROM alpine:3.8

ADD favicon.ico /favicon.ico
RUN mkdir -p /internal/kubernetes/data

ADD static /static
ADD templates /templates
ADD scripts /scripts
ADD kubeyaml-server /kubeyaml
ENTRYPOINT [ "/kubeyaml" ]
