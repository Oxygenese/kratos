module github.com/mars-projects/kratos/contrib/config/kubernetes/v2

go 1.16

require (
	github.com/mars-projects/kratos/v2 v2.1.5
	k8s.io/api v0.23.3
	k8s.io/apimachinery v0.23.3
	k8s.io/client-go v0.23.3
)

replace github.com/mars-projects/kratos/v2 => ../../../
