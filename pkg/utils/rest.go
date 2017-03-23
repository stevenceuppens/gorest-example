package utils

import (
	"fmt"

	restful "github.com/emicklei/go-restful"
)

// ReplyError ...
func ReplyError(code int, message string, err error, res *restful.Response) {
	res.AddHeader("Content-Type", "text/plain")
	res.WriteErrorString(code, message)
	fmt.Println(message, err)
}
