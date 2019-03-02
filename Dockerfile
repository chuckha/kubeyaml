FROM golang:1.12
ADD . /kubeyaml
WORKDIR /kubeyaml
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubeyaml ./cmd/server

FROM alpine:3.8

ADD favicon.ico /favicon.ico
RUN mkdir -p /internal/kubernetes/data
ADD static /static
ADD templates /templates
ADD scripts /scripts
COPY --from=0 /kubeyaml/kubeyaml /kubeyaml
ENTRYPOINT [ "/kubeyaml" ]
