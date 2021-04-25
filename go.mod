module github.com/kbrew-dev/kbrew

go 1.16

replace github.com/graymeta/stow => github.com/kastenhq/stow v0.2.6-kasten.1

require (
	github.com/Masterminds/sprig v2.15.0+incompatible
	github.com/go-git/go-git/v5 v5.2.0
	github.com/google/go-github/v27 v27.0.6
	github.com/kanisterio/kanister v0.0.0-20210224062123-08e898f3dbf3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/openshift/api v0.0.0-20200526144822-34f54f12813a
	github.com/openshift/client-go v0.0.0-20200521150516-05eb9880269c
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.5.4
	k8s.io/api v0.20.4
	k8s.io/client-go v0.20.4
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
