FROM alpine 

RUN set -ex && \
    apk update && \
    apk add ca-certificates

WORKDIR /kube-users

# docker port
EXPOSE 8030

COPY ["k8s-users", "run.sh", "/kube-users/"]

ENTRYPOINT ["/kube-users/run.sh"]
