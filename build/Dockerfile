FROM  ubuntu:14.04

WORKDIR /kube-users

EXPOSE 80

COPY ["k8s-users", "run.sh", "/kube-users/"]

CMD ["/kube-users/run.sh"]
