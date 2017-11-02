# kubencrypt

[![Version Widget]][Version] [![License Widget]][License] [![GoReportCard Widget]][GoReportCard] [![Travis Widget]][Travis] [![DockerHub Widget]][DockerHub]

[Version]: https://github.com/softonic/kubencrypt/releases
[Version Widget]: https://img.shields.io/github/release/softonic/kubencrypt.svg?maxAge=60
[License]: http://www.apache.org/licenses/LICENSE-2.0.txt
[License Widget]: https://img.shields.io/badge/license-APACHE2-1eb0fc.svg
[GoReportCard]: https://goreportcard.com/report/softonic/kubencrypt
[GoReportCard Widget]: https://goreportcard.com/badge/softonic/kubencrypt
[Travis]: https://travis-ci.org/softonic/kubencrypt
[Travis Widget]: https://travis-ci.org/softonic/kubencrypt.svg?branch=master
[DockerHub]: https://hub.docker.com/r/softonic/kubencrypt
[DockerHub Widget]: https://img.shields.io/docker/pulls/softonic/kubencrypt.svg

Letsencrypt on kubernetes.

##### Install

```bash
go get -u github.com/softonic/kubencrypt
```

##### Shell completion

```none
eval "$(kubencrypt --completion-script-${0#-})"
```

##### Out-of-cluster examples

```none
kubencrypt \
--namespace my-namespace \
--ingress my-ingress \
--secret my-secret
```

Same thing with docker:

```none
docker run -it --rm \
-v ~/.kube/config:/root/.kube/config \
-v ~/.config/gcloud:/root/.config/gcloud \
softonic/kubencrypt \
--namespace my-namespace \
--ingress my-ingress \
--secret my-secret
```

##### In-cluster examples

```none
kubectl --namespace my-namespace \
run kubencrypt --restart Never \
--image softonic/kubencrypt \
--image-pull-policy Always -- \
--namespace my-namespace \
--ingress my-ingress \
--secret my-secret
```
