FROM alpine:3.8

ADD templates /templates
ADD static /static
ADD scripts /scripts
ADD favicon.ico /favicon.ico
RUN mkdir -p /internal/kubernetes/data
ADD internal/kubernetes/data/ /internal/kubernetes/data/

ADD kube-validate /kube-validate
ENTRYPOINT [ "/kube-validate" ]
