module bconfig/examples

go 1.25.0

replace github.com/brian-nunez/bconfig => ../

replace github.com/brian-nunez/bconfig/drivers/file => ../drivers/file

replace github.com/brian-nunez/bconfig/drivers/env => ../drivers/env

require (
	github.com/brian-nunez/bconfig v0.0.0
	github.com/brian-nunez/bconfig/drivers/env v0.0.0
	github.com/brian-nunez/bconfig/drivers/file v0.0.0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect
