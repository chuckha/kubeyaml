FROM alpine:3.8

ADD favicon.ico /favicon.ico
RUN mkdir -p /internal/kubernetes/data
ADD internal/kubernetes/data/ /internal/kubernetes/data/

ADD static /static
ADD templates /templates
ADD scripts /scripts
ADD kube-validate /kube-validate
ENTRYPOINT [ "/kube-validate" ]
