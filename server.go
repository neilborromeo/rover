package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

func startServer(frontendFS http.Handler, plan *tfjson.Plan, rso *ResourcesOverview, mapDM *Map, graph Graph) error {
	http.Handle("/", frontendFS)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// simple healthcheck
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"alive": true}`)

	})
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		fileType := strings.Replace(r.URL.Path, "/api/", "", 1)

		var j []byte
		var err error

		enableCors(&w)

		switch fileType {
		case "plan":
			j, err = json.Marshal(plan)
			if err != nil {
				io.WriteString(w, fmt.Sprintf("Error producing plan JSON: %s\n", err))
			}
		case "rso":
			j, err = json.Marshal(rso)
			if err != nil {
				io.WriteString(w, fmt.Sprintf("Error producing rso JSON: %s\n", err))
			}
		case "map":
			j, err = json.Marshal(mapDM)
			if err != nil {
				io.WriteString(w, fmt.Sprintf("Error producing map JSON: %s\n", err))
			}
		case "graph":
			j, err = json.Marshal(graph)
			if err != nil {
				io.WriteString(w, fmt.Sprintf("Error producing graph JSON: %s\n", err))
			}
		default:
			io.WriteString(w, "Please enter a valid file type: plan, rso, map, graph\n")
		}

		w.Header().Set("Content-Type", "application/json")
		io.Copy(w, bytes.NewReader(j))
	})

	log.Println("Rover is running on localhost:9000")

	return http.ListenAndServe(":9000", nil)
}
