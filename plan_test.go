package tweed_test

import (
	"testing"

	"github.com/tweedproject/tweed"
)

func TestPlanInputValidation(t *testing.T) {
	p := tweed.Plan{}
	p.Tweed.Params = map[string]tweed.Parameter{
		"disk": tweed.Parameter{
			Type:    "bytes",
			Minimum: "1Gi",
			Maximum: "100Gi",
		},
		"cpu": tweed.Parameter{
			Type:    "number",
			Minimum: "1",
			Maximum: "4",
		},
	}

	if err := p.ValidateInputs(map[string]string{"disk": "5Gi"}); err != nil {
		t.Errorf("input validation fails: %s", err)
	}
	if err := p.ValidateInputs(map[string]string{"disk": "5Gi", "cpu": "2"}); err != nil {
		t.Errorf("input validation fails: %s", err)
	}

	if err := p.ValidateInputs(map[string]string{"disk": "2Ti"}); err == nil {
		t.Errorf("input validation DOES NOT fail: (no error)")
	}
	if err := p.ValidateInputs(map[string]string{"disk": "2Ti", "cpu": "8"}); err == nil {
		t.Errorf("input validation DOES NOT fail: (no error)")
	}
	if err := p.ValidateInputs(map[string]string{"ram": "80Gi"}); err == nil {
		t.Errorf("input validation DOES NOT fail: (no error)")
	}
}
