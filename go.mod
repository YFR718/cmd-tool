module github.com/YFR718/cmd-tool

go 1.19

require (
	github.com/YFR718/cmd-tool/server/cloud-disk v0.0.0
	github.com/spf13/cobra v1.6.1

)

replace github.com/YFR718/cmd-tool/server/cloud-disk => ./server/cloud-disk

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
