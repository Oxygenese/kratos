module github.com/mars-projects/kratos/contrib/registry/nacos/v2

go 1.16

require (
	github.com/mars-projects/kratos/v2 v2.1.5
	github.com/nacos-group/nacos-sdk-go v1.0.9
)

replace github.com/mars-projects/kratos/v2 => ../../../

replace github.com/buger/jsonparser => github.com/buger/jsonparser v1.1.1
