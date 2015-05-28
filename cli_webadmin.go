package main

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/hashicorp/consul/api"
	"github.com/martini-contrib/render"
)

var s service

func fwrulesWebAdmin(c *cli.Context) {

	s = newService()

	s.cli = c
	s.config = readConfig(c.GlobalString("config"))
	s.client = NewClient(s.config)

	router := martini.Classic()
	router.Get("/status", apiStatus)

	router.NotFound(apiNotFound)

	router.Group("/api", func(api martini.Router) {
		// Machines
		api.Get("/machines", apiMachinesGet)
		// Templates
		api.Get("/templates", apiTemplatesGet)
		api.Post("/templates", apiTemplatesPost)
		api.Put("/templates/**", apiTemplatesPut)
		api.Delete("/templates/**", apiTemplatesDelete)
		// Groups
		api.Get("/groups", apiGroupsGet)
		api.Post("/groups", apiGroupsPost)
		api.Put("/groups/**", apiGroupsPut)
		api.Delete("/groups/**", apiGroupsDelete)
	})

	// Static pages
	router.Use(martini.Static(path.Dir(".")))
	// Renderer
	router.Use(render.Renderer())

	router.RunOnAddr(fmt.Sprintf(":%s", c.String("port")))
}

func apiStatus(r render.Render) {
	r.JSON(200, map[string]bool{"ok": true})
}

func apiMachinesGet(r render.Render) {
	var machines []*Machine

	opts := api.QueryOptions{}

	services, _, err := s.client.catalog.Service("fwrules", "", &opts)
	if err != nil {
		r.JSON(500, []string{})
	}

	for _, service := range services {

		machine, _, err := s.client.Get(pathMachine(service.Node))
		assert(err)

		fmt.Println(machine)

		if err == nil {
			var m Machine
			err = json.Unmarshal(machine, &m)
			if err == nil {
				machines = append(machines, &m)
			}
		}
	}

	r.JSON(200, machines)
}

func apiTemplatesGet() {

}
func apiTemplatesPost() {

}
func apiTemplatesPut() {

}
func apiTemplatesDelete() {

}

func apiGroupsGet() {

}
func apiGroupsPost() {

}
func apiGroupsPut() {

}
func apiGroupsDelete() {

}

func apiNotFound() string {
	return "<html><head><title>Not found!!</title></head><body><p><center style='font-size: 200px; font-family: sans-serif; font-weight: bold;'>404</center></p></body></html>"
}
