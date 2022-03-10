module github.com/mars-projects/kratos/contrib/config/etcd/v2

go 1.16

require (
	github.com/mars-projects/kratos/v2 v2.1.5
	go.etcd.io/etcd/client/v3 v3.5.0
	google.golang.org/grpc v1.44.0
)

replace github.com/mars-projects/kratos/v2 => ../../../
