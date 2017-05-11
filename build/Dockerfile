FROM golang 

WORKDIR /kube-users

# docker port
EXPOSE 8030

COPY ["k8s-users", "run.sh", "/kube-users/"]

CMD ["/kube-users/run.sh"]
