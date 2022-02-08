package main

import (
	"github.com/sustglobal/gost/component"
	"github.com/sustglobal/gost/examples/sample_component"
)

type ServiceConfig struct {
	CustomValue string `env:"CUSTOM_VALUE"`
}

func main() {
	cmp, err := component.NewFromEnv()
	if err != nil {
		panic(err)
	}

	var cfg ServiceConfig
	if err := component.LoadFromEnv(&cfg); err != nil {
		panic(err)
	}

	cmp.HTTPRouter.HandleFunc("/domain", sample_component.DummyHandlerFunc)

	cmp.Run()
}
