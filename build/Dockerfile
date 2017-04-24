FROM golang 

WORKDIR /kube-users

EXPOSE 8030

COPY ["k8s-users", "run.sh", "/kube-users/"]

CMD ["/kube-users/run.sh"]
