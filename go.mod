module github.com/vishal-biyani/kbrew

go 1.15

replace github.com/graymeta/stow => github.com/kastenhq/stow v0.2.6-kasten.1

require (
	github.com/kanisterio/kanister v0.0.0-20210224062123-08e898f3dbf3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/openshift/api v0.0.0-20200526144822-34f54f12813a
	github.com/openshift/client-go v0.0.0-20200521150516-05eb9880269c
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.1
	k8s.io/client-go v0.20.1
)
