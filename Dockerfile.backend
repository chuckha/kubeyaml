FROM golang:1.12
WORKDIR /kubeyaml
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . /kubeyaml
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubeyaml ./cmd/server

FROM alpine:3.8
RUN mkdir -p /internal/kubernetes/data
COPY --from=0 /kubeyaml/kubeyaml /kubeyaml
ENTRYPOINT [ "/kubeyaml" ]
