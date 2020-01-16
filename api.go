package tweed

import (
	"fmt"
	"net/http/httputil"
	"time"

	"github.com/tweedproject/tweed/api"
	"github.com/tweedproject/tweed/random"
	"github.com/tweedproject/tweed/route"
)

func (c *Core) authed(fn func(*route.Request)) func(*route.Request) {
	return func(r *route.Request) {
		if !r.BasicAuthorized(c.HTTPAuthUsername, c.HTTPAuthPassword, c.HTTPAuthRealm) {
			return
		}
		fn(r)
	}
}

func (c *Core) oopsie(r *route.Request, msg string, args ...interface{}) {
	e := fmt.Errorf(msg, args...)
	b, err := httputil.DumpRequest(r.Req, true)
	if err != nil {
		b = []byte(fmt.Sprintf("error parsing http request: %s", err))
	}
	id := random.ID("e")
	if c.oops.index == nil {
		c.oops.index = make([]string, 0)
	}
	if c.oops.faults == nil {
		c.oops.faults = make(map[string]oops)
	}
	c.oops.index = append(c.oops.index, id)
	c.oops.faults[id] = oops{
		handler: r.Req.Method + " " + r.Req.RequestURI,
		remote:  r.Req.RemoteAddr,
		dated:   time.Now(),
		request: b,
		err:     e,
	}

	if len(c.oops.index) > c.oops.max {
		delete(c.oops.faults, c.oops.index[0])
		c.oops.index = c.oops.index[1:]
	}

	r.OK(api.ErrorResponse{
		Err: e.Error(),
		Ref: id,
	})
}

func (c *Core) API() *route.Router {
	r := &route.Router{}

	r.Dispatch("GET /b/status", c.authed(func(r *route.Request) {
		r.OK(map[string]string{"ping": "pong"})
	}))

	r.Dispatch("GET /b/catalog", c.authed(func(r *route.Request) {
		c.UpdateCatalogNumbers()
		r.OK(c.Config.Catalog)
	}))

	r.Dispatch("GET /b/oops/:id", c.authed(func(r *route.Request) {
		e, found := c.oops.faults[r.Args[1]]
		if !found {
			r.Fail(route.NotFound(nil, "error '%s' not found.  this is, itself, an error.  how odd.", r.Args[1]))
			return
		}

		r.OK(api.OopsResponse{
			ID:      r.Args[1],
			Handler: e.handler,
			Remote:  e.remote,
			Request: string(e.request),
			Dated:   e.dated.String(),
			Message: e.message(),
		})
	}))

	r.Dispatch("GET /b/instances", c.authed(func(r *route.Request) {
		l := make([]api.InstanceResponse, 0)
		for _, inst := range c.instances {
			l = append(l, api.InstanceResponse{
				ID:       inst.ID,
				Service:  inst.Plan.Service.ID,
				Plan:     inst.Plan.ID,
				Params:   inst.UserParameters,
				State:    inst.State,
				Log:      inst.Log(),
				Bindings: inst.Bindings,
			})
		}
		r.OK(l)
	}))

	r.Dispatch("GET /b/instances/:id", c.authed(func(r *route.Request) {
		inst, ok := c.instances[r.Args[1]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' not found", r.Args[1]))
			return
		}

		r.OK(api.InstanceResponse{
			ID:       inst.ID,
			Service:  inst.Plan.Service.ID,
			Plan:     inst.Plan.ID,
			Params:   inst.UserParameters,
			State:    inst.State,
			Log:      inst.Log(),
			Bindings: inst.Bindings,
		})
	}))

	r.Dispatch("GET /b/instances/:id/tasks", c.authed(func(r *route.Request) {
		inst, ok := c.instances[r.Args[1]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' not found", r.Args[1]))
			return
		}

		tasks := make([]api.TaskResponse, len(inst.Tasks))
		for i, t := range inst.Tasks {
			tasks[i] = api.TaskResponse{
				Task:     t.id,
				Done:     t.done,
				Exited:   t.exited,
				ExitCode: t.rc,
				Stdout:   t.Stdout(),
				Stderr:   t.Stderr(),
			}
		}

		r.OK(tasks)
	}))

	r.Dispatch("GET /b/instances/:id/tasks/:tid", c.authed(func(r *route.Request) {
		inst, ok := c.instances[r.Args[1]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' not found", r.Args[1]))
			return
		}

		for _, t := range inst.Tasks {
			if t.id == r.Args[2] {
				r.OK(api.TaskResponse{
					Task:     t.id,
					Done:     t.done,
					Exited:   t.exited,
					ExitCode: t.rc,
					Stdout:   t.Stdout(),
					Stderr:   t.Stderr(),
				})
				return
			}
		}

		r.Fail(route.NotFound(nil, "service instance '%s': task '%s' not found", r.Args[1], r.Args[2]))
	}))

	r.Dispatch("GET /b/catalog/:sid/:pid", c.authed(func(r *route.Request) {
		plan, err := c.Config.Catalog.FindPlan(r.Args[1], r.Args[2])
		if err != nil {
			r.Fail(route.Bad(err, "no such service / plan '%s' / '%s'", r.Args[1], r.Args[2]))
			return
		}

		s, err := c.StencilFactory.Get(plan.Tweed.Stencil)
		if err != nil {
			r.Fail(route.Bad(err, err.Error()))
			return
		}

		inst := Instance{
			ID:          "vi-ab-le",
			Plan:        plan,
			Root:        c.Root,
			Prefix:      c.Config.Prefix,
			VaultPrefix: c.Config.Vault.Prefix,
			Stencil:     s,
		}
		if err := inst.Viable(); err != nil {
			r.OK(api.ViabilityResponse{
				Error: fmt.Sprintf("service '%s' / plan '%s' is NOT viable:\n%s\n", plan.Service.Name, plan.Name, err),
			})
			return
		}

		r.OK(api.ViabilityResponse{
			OK: fmt.Sprintf("service '%s' / plan '%s' is viable", plan.Service.Name, plan.Name),
		})
	}))

	r.Dispatch("PUT /b/instances/:id", c.authed(func(r *route.Request) {
		var in api.ProvisionRequest
		if !r.Payload(&in) {
			return
		}

		if err := ValidInstanceID(r.Args[1]); err != nil {
			r.Fail(route.Bad(err, err.Error()))
			return
		}

		plan, err := c.Config.Catalog.FindPlan(in.Service, in.Plan)
		if err != nil {
			r.Fail(route.Bad(err, "no such service / plan '%s' / '%s'", in.Service, in.Plan))
			return
		}

		s, err := c.StencilFactory.Get(plan.Tweed.Stencil)
		if err != nil {
			r.Fail(route.Bad(err, err.Error()))
			return
		}

		ref, err := c.Provision(&Instance{
			ID:             r.Args[1],
			Plan:           plan,
			Root:           c.Root,
			Prefix:         c.Config.Prefix,
			VaultPrefix:    c.Config.Vault.Prefix,
			UserParameters: in.Params,
			Stencil:        s,
		})
		if err != nil {
			c.oopsie(r, "unable to provision a %s / %s service instance: %s", in.Service, in.Plan, err)
			return
		}

		r.OK(api.ProvisionResponse{
			OK:  "service instance scheduled for provisioning; thank you for your patience.",
			Ref: ref,
		})
	}))

	r.Dispatch("GET /b/instances/:id/files", c.authed(func(r *route.Request) {
		inst, ok := c.instances[r.Args[1]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' not found", r.Args[1]))
			return
		}

		files, err := inst.Files()
		if err != nil {
			c.oopsie(r, "unable to retrieve files: %s", err)
			return
		}
		r.OK(files)
	}))

	r.Dispatch("PUT /b/instances/:id/state/:state", c.authed(func(r *route.Request) {
		inst, ok := c.instances[r.Args[1]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' not found", r.Args[1]))
			return
		}

		was := inst.State
		inst.State = r.Args[2]
		r.Success("state manually changed from '%s' to '%s'", was, inst.State)
	}))

	r.Dispatch("DELETE /b/instances/:id", c.authed(func(r *route.Request) {
		ref, gone, err := c.Deprovision(r.Args[1])
		if err != nil {
			c.oopsie(r, "unable to deprovision service instance %s", r.Args[1])
			return
		}

		res := api.DeprovisionResponse{}
		if gone {
			res.Error = fmt.Sprintf("service instance '%s' already deprovisioned", r.Args[1])
			res.Gone = true

		} else {
			res.OK = "service instance scheduled for deprovisioning; thank you for your patience."
			res.Ref = ref
		}
		r.OK(res)
	}))

	r.Dispatch("DELETE /b/instances/:id/log", c.authed(func(r *route.Request) {
		inst, ok := c.instances[r.Args[1]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' not found", r.Args[1]))
			return
		}

		if err := inst.Purge(); err != nil {
			c.oopsie(r, "unable to purge service instance %s: %s", inst.ID, err)
			return
		}

		delete(c.instances, inst.ID)
		r.Success("service instance '%s' purged", inst.ID)
	}))

	r.Dispatch("GET /b/instances/:id/bindings", c.authed(func(r *route.Request) {
		inst, ok := c.instances[r.Args[1]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' not found", r.Args[1]))
			return
		}

		r.OK(api.BindingsResponse{
			Bindings: inst.Bindings,
		})
	}))

	r.Dispatch("PUT /b/instances/:id/bindings/:bid", c.authed(func(r *route.Request) {
		ref, err := c.Bind(r.Args[1], r.Args[2])
		if err != nil {
			c.oopsie(r, "unable to bind service instance '%s': %s", r.Args[1], err)
			return
		}

		r.OK(api.BindResponse{
			OK:  "service instance bind operation scheduled; thank you for your patience.",
			Ref: ref,
		})
	}))

	r.Dispatch("GET /b/instances/:id/bindings/:bid", c.authed(func(r *route.Request) {
		inst, ok := c.instances[r.Args[1]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' not found", r.Args[1]))
			return
		}

		binding, ok := inst.Bindings[r.Args[2]]
		if !ok {
			r.Fail(route.NotFound(nil, "service instance '%s' binding '%s' not found", r.Args[1], r.Args[2]))
			return
		}
		r.OK(api.BindingResponse{
			Binding: binding,
		})
	}))

	r.Dispatch("DELETE /b/instances/:id/bindings/:bid", c.authed(func(r *route.Request) {
		ref, err := c.Unbind(r.Args[1], r.Args[2])
		if err != nil {
			c.oopsie(r, "unable to unbind service instance '%s': %s", r.Args[1], err)
			return
		}

		r.OK(api.UnbindResponse{
			OK:  "service instance unbind operation scheduled; thank you for your patience.",
			Ref: ref,
		})
	}))

	return r
}
