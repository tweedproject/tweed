package tweed

type Service struct {
	Name        string   `json:"name"`
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`

	Bindable             bool `json:"bindable"`
	InstancesRetrievable bool `json:"instances_retrievable"`
	BindingsRetrievable  bool `json:"bindings_retrievable"`
	AllowContextUpdates  bool `json:"allow_context_updates"`
	PlanUpdateable       bool `json:"plan_updateable"`

	Metadata map[string]interface{} `json:"metadata"`

	Plans []*Plan `json:"plans"`
}
