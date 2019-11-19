package tweed_test

import (
	"testing"

	"github.com/tweedproject/tweed"
)

func TestTextParameters(t *testing.T) {
	p := tweed.Parameter{
		Type: tweed.TextParameter,
	}

	if err := p.Wellformed(); err != nil {
		t.Fatalf("text parameter (no validation) is not considered well-formed: %s", err)
	}
}

func TestBytesParameters(t *testing.T) {
	tests := []struct {
		min string
		max string
		ok  bool
	}{
		{min: "1Ki", max: "1Mi", ok: true},
		{min: "1Gi", max: "1Ti", ok: true},
		{min: "1KiB", max: "1MiB", ok: true},
		{min: "1GiB", max: "1TiB", ok: true},
		{min: "1KB", max: "1MB", ok: true},
		{min: "1GB", max: "1TB", ok: true},
		{min: "1K", max: "1M", ok: true},
		{min: "1G", max: "1T", ok: true},
		{min: "1ki", max: "1mib", ok: true},
		{min: "1gb", max: "1t", ok: true},
		{min: "", max: "5Gi", ok: true},
		{min: "5Ki", max: "", ok: true},
		{min: "1024", max: "2048", ok: true},

		{min: "5 bytes", max: "", ok: false},

		{min: "1Mi", max: "1ki", ok: false},
	}

	for _, test := range tests {
		p := tweed.Parameter{
			Type:    tweed.BytesParameter,
			Minimum: test.min,
			Maximum: test.max,
		}

		err := p.Wellformed()
		if test.ok && err != nil {
			t.Errorf("bytes parameter (%s ... %s) is mistakenly considered NOT well-formed: %s", test.min, test.max, err)
		} else if !test.ok && err == nil {
			t.Errorf("bytes parameter (%s ... %s) is MISTAKENLY considered well-formed: (no error)", test.min, test.max)
		}
	}
}

func TestNumericParameters(t *testing.T) {
	tests := []struct {
		min string
		max string
		ok  bool
	}{
		{min: "50", max: "100", ok: true},
		{min: "", max: "7", ok: true},
		{min: "4", max: "", ok: true},

		{min: "one", max: "three", ok: false},

		{min: "100", max: "50", ok: false},
	}

	for _, test := range tests {
		p := tweed.Parameter{
			Type:    tweed.BytesParameter,
			Minimum: test.min,
			Maximum: test.max,
		}

		err := p.Wellformed()
		if test.ok && err != nil {
			t.Errorf("numeric parameter (%s ... %s) is mistakenly considered NOT well-formed: %s", test.min, test.max, err)
		} else if !test.ok && err == nil {
			t.Errorf("numeric parameter (%s ... %s) is MISTAKENLY considered well-formed: (no error)", test.min, test.max)
		}
	}
}

func TestDurationParameters(t *testing.T) {
	tests := []struct {
		min string
		max string
		ok  bool
	}{
		{min: "1s", max: "1m", ok: true},
		{min: "1h", max: "1d", ok: true},
		{min: "1w", max: "1y", ok: true},
		{min: "1S", max: "1M", ok: true},
		{min: "1H", max: "1D", ok: true},
		{min: "1W", max: "1Y", ok: true},
		{min: "", max: "100d", ok: true},
		{min: "100d", max: "", ok: true},
		{min: "40", max: "50", ok: true},

		{min: "5 days", max: "", ok: false},

		{min: "8h", max: "1h", ok: false},
	}

	for _, test := range tests {
		p := tweed.Parameter{
			Type:    tweed.DurationParameter,
			Minimum: test.min,
			Maximum: test.max,
		}

		err := p.Wellformed()
		if test.ok && err != nil {
			t.Errorf("duration parameter (%s ... %s) is mistakenly considered NOT well-formed: %s", test.min, test.max, err)
		} else if !test.ok && err == nil {
			t.Errorf("duration parameter (%s ... %s) is MISTAKENLY considered well-formed: (no error)", test.min, test.max)
		}
	}
}
