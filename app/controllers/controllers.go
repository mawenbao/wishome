package controllers

import (
    "net"
    "strings"
    "github.com/robfig/revel"
)

func init() {
    // load interceptors
    revel.OnAppStart(LoadInterceptors)
}

func GetRemoteAddr(c *revel.Controller) string {
    // check X-Forwarded-For header
    if forwds := c.Request.Header.Get("X-Forwarded-For"); "" == forwds {
        // split host and port
        host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
        if nil != err {
            revel.ERROR.Printf("failed to split host:port string %s: %s", c.Request.RemoteAddr, err)
            return strings.TrimSpace(c.Request.RemoteAddr)
        }
        return strings.TrimSpace(host)
    } else {
        // pick the first elem in X-Forwarded-For header separated by a comma(,)
        return strings.TrimSpace(strings.Split(forwds, ",")[0])
    }
}

