package main

import (
	"os"
	"time"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type patience struct {
	instance string
	task     string
	until    string
	negate   bool
	nap      int
	slept    int
	max      int
	quiet    bool
}

func (p patience) printf(f string, args ...interface{}) {
	if !p.quiet && !Tweed.JSON {
		fmt.Printf(f, args...)
	}
}

func (p *patience) sleep() {
	if p.nap == 0 {
		p.nap = 2
	}

	if p.task == "" { // instance wait
		if p.negate {
			p.printf("instance @W{%s} is still @W{%s}; sleeping for another %d seconds...", p.instance, p.until, p.nap)
		} else {
			p.printf("instance @W{%s} is not yet @W{%s}; sleeping for another %d seconds...", p.instance, p.until, p.nap)
		}

	} else { // task wait
		p.printf("task @W{%s}/@W{%s} still running; sleeping for another %d seconds...", p.instance, p.task, p.nap)
	}

	if p.max > 0 {
		p.printf(" (%ds left until we give up)", p.max-p.slept)
	}
	p.printf("\n")

	time.Sleep(time.Duration(p.nap) * time.Second)
	p.slept += p.nap
}

func awaitfn(p patience, fn func() bool) bool {
	p.slept = 0
	for {
		if done := fn(); done {
			return true
		}

		if p.max > 0 && p.max <= p.slept {
			p.printf("exceeded maximum wait time of @Y{%d seconds}.\n", p.max)
			return false
		}
		p.sleep()
	}
}

func await(c *client, p patience) bool {
	if p.task == "" { // instance wait
		return awaitfn(p, func() bool {
			var out api.InstanceResponse
			c.GET("/b/instances/"+p.instance, &out)
			if p.negate && out.State != p.until {
				p.printf("instance @C{%s} is no longer @G{%s}.\n", p.instance, p.until)
				return true
			}
			if !p.negate && out.State == p.until {
				p.printf("instance @C{%s} is now @G{%s}.\n", p.instance, p.until)
				return true
			}

			return false
		})

	} else { // task wait
		return awaitfn(p, func() bool {
			var out api.TaskResponse
			c.GET("/b/instances/"+p.instance+"/tasks/"+p.task, &out)
			if out.ExitCode != 0 {
				p.printf("@r{task %s/%s failed:}\n%s\n",
					p.instance, p.task, out.Stderr)
				os.Exit(out.ExitCode)
			}

			return out.Done
		})
	}

	panic("maflformed await() / patience call")
}
