package results

import (
    "net/http"
    "github.com/robfig/revel"
)

type ImagePngResult []byte

func (r ImagePngResult) Apply(req *revel.Request, resp *revel.Response) {
    resp.WriteHeader(http.StatusOK, "image/png")
    resp.Out.Write(r)
}

