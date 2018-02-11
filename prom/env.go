package prom

import "os"

var (
	ZTOHostName string = ""
)

func init() {
	ZTOHostName = os.Getenv("ZTO_HOST_NAME")
	if ZTOHostName == "" {
		panic("please set environment variable ZTO_HOST_NAME")
	}
	logger.Infof("ZTO_HOST_NAME:\t%s", ZTOHostName)
}
