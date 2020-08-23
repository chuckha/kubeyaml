FROM golang:1.15
WORKDIR /kubeyaml
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . /kubeyaml
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubeyaml ./cmd/server

FROM alpine:3.9
ARG USER=default
ENV HOME /home/$USER

# install sudo as root
RUN apk add --update sudo

# add new user
RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

USER $USER
WORKDIR $HOME

COPY --from=0 /kubeyaml/kubeyaml ./kubeyaml
ENTRYPOINT [ "./kubeyaml" ]