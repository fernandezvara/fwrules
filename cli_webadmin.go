package main

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"
	"github.com/hashicorp/consul/api"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

var s service

func fwrulesWebAdmin(c *cli.Context) {

	s.cli = c
	s.config = readConfig(c.GlobalString("config"))
	s.client = NewClient(s.config)

	router := martini.Classic()
	router.Get("/status", apiStatus)

	router.NotFound(apiNotFound)

	router.Group("/api", func(api martini.Router) {
		// Machines
		api.Get("/:fwid/machines", apiMachinesGet)
		// RuleSets
		api.Get("/:fwid/rulesets", apiRuleSetsGetAll)
		api.Get("/:fwid/rulesets/:id", apiRuleSetsGet)
		api.Post("/:fwid/rulesets", binding.Json(RuleSet{}), apiRuleSetsPost)
		api.Put("/:fwid/rulesets/:id", binding.Json(RuleSet{}), apiRuleSetsPut)
		api.Delete("/:fwid/rulesets/**", apiRuleSetsDelete)
	})

	// Static pages
	router.Use(martini.Static(path.Dir(".")))
	// Renderer
	router.Use(render.Renderer())

	router.RunOnAddr(fmt.Sprintf("%s:%s", c.String("ip"), c.String("port")))
}

func apiStatus(r render.Render) {
	r.JSON(200, map[string]bool{"ok": true})
}

func apiMachinesGet(r render.Render, params martini.Params) {
	var machines []*Machine

	opts := api.QueryOptions{}

	services, _, err := s.client.catalog.Service(fmt.Sprintf("fwrules-%s", params["fwid"]), "", &opts)
	if err != nil {
		r.JSON(500, []string{})
	}

	for _, service := range services {

		machine, _, err := s.client.Get(pathMachine(params["fwid"], service.Node))
		assert(err)

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

func apiRuleSetsGetAll(r render.Render, params martini.Params) {
	var rulesets []RuleSet
	entries, _, err := s.client.List(pathRuleSet(params["fwid"], ""))
	for _, entry := range entries {
		var t RuleSet
		err = json.Unmarshal(entry.Value, &t)
		if err == nil {
			rulesets = append(rulesets, t)
		}
	}
	r.JSON(200, rulesets)
}

func apiRuleSetsGet(r render.Render, params martini.Params) {
	entry, exists, err := s.client.Get(pathRuleSet(params["fwid"], params["id"]))
	if err != nil {
		r.JSON(500, map[string]string{"error": "cannot get data"})
		return
	}
	if exists == false {
		r.JSON(404, map[string]string{"error": "not found"})
		return
	}
	var t RuleSet
	err = json.Unmarshal(entry, &t)
	if err != nil {
		r.JSON(500, map[string]string{"error": "data returned cannot be unmarshalled"})
		return
	}
	r.JSON(200, t)
}

func apiRuleSetsPost(ruleset RuleSet, r render.Render, params martini.Params) {
	if ruleset.Name == "" {
		r.JSON(400, map[string]string{"error": "name cannot be blank"})
		return
	}
	_, exists, err := s.client.Get(pathRuleSet(params["fwid"], ruleset.Name))
	assert(err)
	if exists == true {
		r.JSON(409, map[string]string{"error": "conflict"})
		return
	}
	data, err := json.Marshal(ruleset)
	if err != nil {
		r.JSON(500, map[string]string{"error": "cannot marshal data"})
		return
	}
	err = s.client.Set(pathRuleSet(params["fwid"], ruleset.Name), data)
	if err != nil {
		r.JSON(500, map[string]string{"error": "cannot set data"})
		return
	}
	r.JSON(201, ruleset)
}

func apiRuleSetsPut(ruleset RuleSet, r render.Render, params martini.Params) {
	if params["id"] != ruleset.Name {
		r.JSON(409, map[string]string{"error": "conflict"})
	}
	data, err := json.Marshal(ruleset)
	if err != nil {
		r.JSON(500, map[string]string{"error": "cannot marshal data"})
		return
	}
	err = s.client.Set(pathRuleSet(params["fwid"], ruleset.Name), data)
	if err != nil {
		r.JSON(500, map[string]string{"error": "cannot set data"})
		return
	}
	r.JSON(204, nil)
}

func apiRuleSetsDelete(ruleset RuleSet, r render.Render, params martini.Params) {
	err := s.client.Delete(pathRuleSet(params["fwid"], ruleset.Name))
	if err != nil {
		r.JSON(500, map[string]string{"error": "could not delete"})
	}
	r.JSON(200, nil)
}

func apiNotFound() string {
	return "<html><head><title>Not found!!</title></head><body><p><center style='font-size: 200px; font-family: sans-serif; font-weight: bold;'>404</center></p></body></html>"
}
