# k8s-users

## 背景
在适用ABAC的认证方式，为Kubernetes集群创建用户时，通常需要经过以下几个步骤：
- 管理根据用户信息为用户分配一个用户名，根据Kubernetes集群的CA证书生成专属这个用户的tls证书
- 如果需要为此用户创建一个新namespace，还需要使用`kubectl create namespace`为用户创建一个namespace
- 修改Kubernetes master节点上policy文件，配置用户的权限
- 重启apiserver进程
- 将用户名以及tls证书通过邮件等方式发给用户

## 特性
此工具可以将这个流程变成自动化的形式，管理员只需要向接口发送一个HTTP的请求即可完成账户开通流程，具体特性如下：

- HTTP接口：工具会对外提供一个HTTP的接口供管理员调用，接口中定义了用户名，所属namespace，邮箱地址等信息：

  ```bash
  cur -XPOST http://kube-master:<port>/users -d '
  {
    "username":"admin",
    "namespace":"admin",
    "email":"admin@domain.com"
  }'
  ```

  - usernam: **required** , 需要管理员提供用户名信息。
  - namespace: **required** , 需要管理员提供所属的namespace，如果此namespace不存在将会重新创建一个，也可以使用*号表示不限制namespace。
  - email: **required**, 账号开通后将会向此邮箱地址发送一封欢迎邮件，并且在附件中包含tls key。

- 发送邮件：账号创建后，通过`mail`命令，将欢迎信息以及tls信息发送给用户。
- 部署：工具将会以[addon](https://github.com/kubernetes/kubernetes/tree/master/cluster/addons)的形式部署在Kubernetes集群的master节点，将集群的`ca.pem`和`ca-key.pem`以hostpath的方式mount到Pod中。


## TODO

- 工具的接口需要通过认证才可以正常的接收请求
- 提供方便操作的web界面
