package api

import (
	"fmt"
)

type ErrorResponse struct {
	Err string `json:"error"`
	Ref string `json:"ref"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s (ref %s)", e.Err, e.Ref)
}

type ProvisionResponse struct {
	OK    string `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
	Ref   string `json:"ref"`
}

type DeprovisionResponse struct {
	OK    string `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
	Gone  bool   `json:"gone"`
	Ref   string `json:"ref"`
}

type PurgeResponse struct {
	OK    string `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
}

type BindResponse struct {
	OK    string `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
	Ref   string `json:"ref"`
}

type UnbindResponse struct {
	OK    string `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
	Ref   string `json:"ref"`
}

type TaskResponse struct {
	Task     string `json:"task"`
	Done     bool   `json:"done"`
	Exited   bool   `json:"exited"`
	ExitCode int    `json:"exit_code"`
	Ref      string `json:"ref"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

type InstanceResponse struct {
	ID      string                 `json:"id"`
	Service string                 `json:"service"`
	Plan    string                 `json:"plan"`
	Params  map[string]interface{} `json:"params"`

	State string `json:"state"`
	Log   string `json:"log"`

	Bindings map[string]map[string]interface{} `json:"bindings"`

	Tasks []TaskResponse `json:"tasks"`
}

type BindingsResponse struct {
	Bindings map[string]map[string]interface{} `json:"bindings"`
}

type BindingResponse struct {
	Binding map[string]interface{} `json:"binding"`
}

type OopsResponse struct {
	ID      string `json:"id"`
	Handler string `json:"handler"`
	Remote  string `json:"remote"`
	Request string `json:"request"`
	Dated   string `json:"dated"`
	Message string `json:"message"`
}

type FileResponse struct {
	Filename    string `json:"filename"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Contents    string `json:"contents"`
}

type ViabilityResponse struct {
	OK    string `json:"ok,omitempty"`
	Error string `json:"error,omitempty"`
}
