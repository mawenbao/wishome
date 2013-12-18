package template

import (
    "bytes"
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

func LoadTempate(name, templateContent string, args interface{}) []byte {
    buff := new(bytes.Buffer)
    t := template.New(name)

    t, err := t.Parse(templateContent)
    if nil != err {
        revel.ERROR.Printf("failed to parse template %s, template content[%s]", err, templateContent)
        return nil
    }

    err = t.Execute(buff, args)
    if nil != err {
        revel.ERROR.Printf("failed to execute template %s, template content %s", err, templateContent)
        return nil
    }

    return buff.Bytes()
}

