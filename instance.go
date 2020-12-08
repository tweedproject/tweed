package tweed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// don't forget to update the api/responses.go copies as well!
	QuietState          = "quiet"
	ProvisioningState   = "provisioning"
	DeprovisioningState = "deprovisioning"
	GoneState           = "gone"
	BindingState        = "binding"
	UnbindingState      = "unbinding"
)

type Instance struct {
	ID    string
	Plan  *Plan
	State string

	Root        string
	Prefix      string
	VaultPrefix string

	UserParameters map[string]interface{}
	Bindings       map[string]map[string]interface{}

	Tasks []*task
}

type instancemf struct {
	Tweed struct {
		Prefix    string `json:"prefix"`
		Instance  string `json:"instance"`
		Service   string `json:"service"`
		ServiceID string `json:"service_id"`
		Plan      string `json:"plan"`
		PlanID    string `json:"plan_id"`
		Vault     string `json:"vault"`

		Ops  map[string]interface{} `json:"ops"`
		User map[string]interface{} `json:"user"`
	} `json:"tweed"`
}

func (i *Instance) path(rel string) string {
	return fmt.Sprintf("%s/%s", i.Root, rel)
}

func (i *Instance) my(rel string) string {
	if rel != "" {
		rel = "/" + strings.TrimPrefix(rel, "/")
	}
	return i.path("data/instances/" + i.ID + rel)
}

func (i *Instance) env(env []string) []string {
	env = append(env, "HOME="+i.path(""))
	env = append(env, "PATH="+os.Getenv("PATH"))
	env = append(env, "LANG="+os.Getenv("LANG"))
	env = append(env, "INFRASTRUCTURE="+i.path("etc/infrastructures/"+i.Plan.Tweed.Infrastructure))
	env = append(env, "STENCIL="+i.path("etc/stencils/"+i.Plan.Tweed.Stencil))
	env = append(env, "WORKSPACE="+i.my(""))
	env = append(env, "VAULT="+i.VaultPrefix+"/"+i.ID)
	env = append(env, "INPUTS=instance.mf")
	return env
}

func ParseInstance(cat Catalog, root string, b []byte) (Instance, error) {
	var in instancemf

	err := json.Unmarshal(b, &in)
	if err != nil {
		return Instance{}, err
	}

	p, err := cat.FindPlan(in.Tweed.ServiceID, in.Tweed.PlanID)
	if err != nil {
		return Instance{}, err
	}

	inst := Instance{
		ID:             in.Tweed.Instance,
		Root:           root,
		Plan:           p,
		UserParameters: in.Tweed.User,
		State:          QuietState,
	}
	b, err = ioutil.ReadFile(inst.my("lifecycle/data/state"))
	if err == nil {
		inst.State = strings.TrimSpace(string(b))
	}

	return inst, nil
}

func (i *Instance) lookupBindings(id string) error {
	if i.Bindings == nil {
		i.Bindings = make(map[string]map[string]interface{})
	}

	b, err := run1(Exec{
		Run: i.path("bin/bindings"),
		Env: i.env([]string{"BINDING=" + id}),
	})
	if err != nil {
		return err
	}

	var all map[string]map[string]interface{}
	err = json.Unmarshal(b, &all)
	if err != nil {
		return err
	}

	for _, bindings := range all {
		for id, raw := range bindings {
			s, ok := raw.(string)
			if !ok {
				return fmt.Errorf("binding %s/%s is not a string")
			}
			var v map[string]interface{}
			if err := json.Unmarshal([]byte(s), &v); err != nil {
				return err
			}
			i.Bindings[id] = v
		}
	}
	return nil
}

func (i *Instance) LookupBindings() error {
	return i.lookupBindings("")
}

func (i *Instance) LookupBinding(id string) error {
	return i.lookupBindings(id)
}

func (i *Instance) Log() string {
	b, _ := ioutil.ReadFile(i.my("log"))
	return string(b)
}

func (i *Instance) do(cmd, begin, middle, end string) (*task, error) {
	if begin != "" && i.State != begin {
		return nil, fmt.Errorf("service instance '%s' is currently %s", i.ID, i.State)
	}

	i.State = middle
	t := background(Exec{
		Run: i.path(cmd),
		Env: i.env(nil),
	}, func() {
		fmt.Printf("updating state to '%s'\n", end)
		i.State = end
	})

	i.Tasks = append(i.Tasks, t)
	return t, nil
}
func (i *Instance) Provision() (*task, error) {
	if err := i.Viable(); err != nil {
		return nil, err
	}

	var out instancemf

	out.Tweed.Prefix = i.Prefix
	out.Tweed.Instance = i.ID
	out.Tweed.Service = i.Plan.Service.Name
	out.Tweed.ServiceID = i.Plan.Service.ID
	out.Tweed.Plan = i.Plan.Name
	out.Tweed.PlanID = i.Plan.ID
	out.Tweed.Vault = `(( concat "` + i.VaultPrefix + `/" tweed.instance ))`
	out.Tweed.Ops = i.Plan.Tweed.Config
	out.Tweed.User = i.UserParameters

	input, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}

	root := i.my("")
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(root+"/instance.mf", input, 0666); err != nil {
		return nil, err
	}

	return i.do("bin/provision", "", ProvisioningState, QuietState)
}

func (i *Instance) Bind(id string) (*task, error) {
	if i.State != QuietState {
		return nil, fmt.Errorf("service instance '%s' is currently %s", i.ID, i.State)
	}

	if err := i.Viable(); err != nil {
		return nil, err
	}

	i.State = BindingState
	t := background(Exec{
		Run: i.path("bin/bind"),
		Env: i.env([]string{
			"BINDING=" + id,
			"OVERRIDES=" + i.CredentialOverrides(),
		}),
	}, func() {
		i.State = QuietState
		if err := i.LookupBinding(id); err != nil {
			fmt.Fprintf(os.Stderr, "failed to look up newly-created binding %s/%s: %s", i.ID, id, err)
		}
	})

	i.Tasks = append(i.Tasks, t)
	return t, nil
}

func (i *Instance) Unbind(id string) (*task, error) {
	if i.State != QuietState {
		return nil, fmt.Errorf("service instance '%s' is currently %s", i.ID, i.State)
	}

	if err := i.Viable(); err != nil {
		return nil, err
	}

	i.State = UnbindingState
	t := background(Exec{
		Run: i.path("bin/unbind"),
		Env: i.env([]string{"BINDING=" + id}),
	}, func() {
		i.State = QuietState
		delete(i.Bindings, id)
	})

	i.Tasks = append(i.Tasks, t)
	return t, nil
}

func (i *Instance) Deprovision() (*task, error) {
	if err := i.Viable(); err != nil {
		return nil, err
	}

	return i.do("bin/deprovision", QuietState, DeprovisioningState, GoneState)
}

func (i *Instance) Purge() error {
	if i.State != GoneState {
		return fmt.Errorf("service instance '%s' is currently %s", i.ID, i.State)
	}

	return os.RemoveAll(i.my(""))
}

func (i *Instance) Viable() error {
	out, err := run1(Exec{
		Run: i.path("bin/viable"),
		Env: i.env(nil),
	})
	if err != nil {
		return fmt.Errorf("stencil viability check failed: %s", string(out))
	}
	return nil
}

func (i *Instance) CredentialOverrides() string {
	if i.Plan.Tweed.Credentials == nil {
		return `{}`
	}
	out := map[string]interface{}{
		"credentials": i.Plan.Tweed.Credentials,
	}
	if b, err := json.Marshal(&out); err != nil {
		return `{}`
	} else {
		return string(b)
	}
}

func (i *Instance) Files() ([]File, error) {
	out, err := run1(Exec{
		Run: i.path("bin/files"),
		Env: i.env(nil),
	})
	if err != nil {
		return nil, err
	}

	var f struct {
		Files []File `json:"files"`
	}
	return f.Files, json.Unmarshal(out, &f)
}

func (i Instance) IsBusy() bool {
	return i.State == ProvisioningState || i.State == DeprovisioningState || i.State == BindingState || i.State == UnbindingState
}

func (i Instance) IsQuiet() bool {
	return i.State == QuietState
}

func (i Instance) IsGone() bool {
	return i.State == GoneState
}
