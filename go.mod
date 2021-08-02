module github.com/kbrew-dev/kbrew

go 1.16

replace github.com/graymeta/stow => github.com/kastenhq/stow v0.2.6-kasten.1

require (
	cloud.google.com/go/bigquery v1.8.0 // indirect
	github.com/BurntSushi/toml v0.3.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/briandowns/spinner v1.16.0
	github.com/go-git/go-git/v5 v5.2.0
	github.com/google/go-cmp v0.5.4
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/kanisterio/kanister v0.0.0-20210224062123-08e898f3dbf3
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/mikefarah/yq/v4 v4.9.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/openshift/api v0.0.0-20200526144822-34f54f12813a
	github.com/openshift/client-go v0.0.0-20200521150516-05eb9880269c
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.6.1
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
