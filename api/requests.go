package api

type ProvisionRequest struct {
	Service string                 `json:"service"`
	Plan    string                 `json:"plan"`
	Params  map[string]interface{} `json:"params"`
}
