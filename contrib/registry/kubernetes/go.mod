module github.com/mars-projects/kratos/contrib/registry/kubernetes/v2

go 1.16

require (
	github.com/json-iterator/go v1.1.12
	github.com/mars-projects/kratos/v2 v2.0.0-00010101000000-000000000000
	k8s.io/api v0.23.1
	k8s.io/apimachinery v0.23.1
	k8s.io/client-go v0.23.1
)

replace github.com/mars-projects/kratos/v2 => ../../../
