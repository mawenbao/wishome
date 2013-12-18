package template

import (
    "html/template"
    "github.com/robfig/revel"
)

func init() {
    // init custom template functions
    revel.TemplateFuncs["unescape"] = unescape
}

func unescape (x string) interface{} {
    return template.HTML(x)
}

