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
    // check first element of X-Forwarded-For header
    if forwds := c.Request.Header.Get("X-Forwarded-For"); "" == forwds {
        // check X-Real-Ip
        if realIp := c.Request.Header.Get("X-Real-Ip"); "" != realIp {
            return strings.TrimSpace(realIp)
        }
        // use Request.RemoteAddr, split host and port
        host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
        if nil != err {
            revel.ERROR.Printf("failed to split host:port string %s: %s", c.Request.RemoteAddr, err)
            return strings.TrimSpace(c.Request.RemoteAddr)
        }
        return strings.TrimSpace(host)
    } else {
        return strings.TrimSpace(forwds)
    }
}

