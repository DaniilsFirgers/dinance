package cron

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type Cron struct {
	cron *cron.Cron
}

func New() *Cron {
	return &Cron{
		cron: cron.New(cron.WithSeconds(), cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		)),
	}
}

func (c *Cron) Start() {
	c.cron.Start()
}

func (c *Cron) Stop() context.Context {
	return c.cron.Stop()
}

func (c *Cron) AddFunc(name string, spec string, cmd func()) (cron.EntryID, error) {
	wrappedCmd := func() {
		start := time.Now()
		log.Printf("%s started at %s", name, start.Format(time.RFC3339))
		cmd()
		duration := time.Since(start)
		log.Printf("%s finished in %s", name, duration)
	}

	return c.cron.AddFunc(spec, wrappedCmd)
}
