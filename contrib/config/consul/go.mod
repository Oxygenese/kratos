module github.com/mars-projects/kratos/contrib/config/consul/v2

go 1.15

require (
	github.com/hashicorp/consul/api v1.10.0
	github.com/mars-projects/kratos/v2 v2.1.5
)

replace github.com/mars-projects/kratos/v2 => ../../../
