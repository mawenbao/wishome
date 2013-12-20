package controllers

import (
    "github.com/robfig/revel"
)

func init() {
    // load interceptors
    revel.OnAppStart(LoadInterceptors)
}

