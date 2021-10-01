package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Basic structure that reflects some information of the Node structure.
// This is constructed and passed with go templating, so that the WebUI shows some
// information about the node.
type NodeInfo struct {
	Data       map[string]string
	Neighbours []string
	Time       string
}

func (N *Node) webHandler(w http.ResponseWriter, r *http.Request) {
	information := NodeInfo{
		Data:       N.Data,
		Neighbours: N.Neighbours,
		Time:       time.Now().Format(time.RFC1123),
	}
	tpl := template.Must(template.ParseGlob("web/*"))
	if err := tpl.ExecuteTemplate(w, "node_information.html", information); err != nil {
		log.Fatalf("Could not start web server: %v\n", err)
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Webserver couldn't parse form: %v\n", err)
	} else if len(r.Form["form_key"]) != 0 && len(r.Form["form_value"]) != 0 {
		err := N.Send(fmt.Sprintf("tcp://127.0.0.1:%d", N.Port), Rumour{
			RequestType: UpdateData,
			Key:         r.Form["form_key"][0],
			NewValue:    r.Form["form_value"][0],
		})
        if err != nil {
            log.Printf("could not send values, err: %v\n", err)
        }
	}
}
