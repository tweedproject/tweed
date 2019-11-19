package tweed

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	TextParameter     = "text"
	BytesParameter    = "bytes"
	NumericParameter  = "number"
	DurationParameter = "duration"
)

type MalformedParameterSpec struct {
	msg string
}

func (e MalformedParameterSpec) Error() string {
	return fmt.Sprintf("malformed parameter specification: %s", e.msg)
}

type WrongParameterValueType struct {
	expect string
	got    string
}

func (e WrongParameterValueType) Error() string {
	return fmt.Sprintf("incorrect parameter value type; expected %s, but got %s", e.expect, e.got)
}

type ParameterValueOutOfRange struct {
	value string
	msg   string
}

func (e ParameterValueOutOfRange) Error() string {
	return fmt.Sprintf("parameter value '%s' is out of range: %s", e.value, e.msg)
}

type Parameter struct {
	Type    string `json:"type"`
	Minimum string `json:"minimum"`
	Maximum string `json:"maximum"`
}

func (p Parameter) Wellformed() error {
	switch p.Type {
	case TextParameter:
		return nil

	case BytesParameter:
		min := p.Minimum == "" || p.parsesAsBytes(p.Minimum)
		max := p.Maximum == "" || p.parsesAsBytes(p.Maximum)

		if !min && !max {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid minimum / maximum specified: '%s' / '%s' are not sizes", p.Minimum, p.Maximum),
			}
		}
		if !min {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid minimum specified: '%s' is not a size", p.Minimum),
			}
		}
		if !min {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid maximum specified: '%s' is not a size", p.Maximum),
			}
		}
		if p.Minimum != "" && p.Maximum != "" {
			floor, _ := p.parseAsBytes(p.Minimum)
			ceil, _ := p.parseAsBytes(p.Maximum)

			if ceil < floor {
				return MalformedParameterSpec{
					msg: fmt.Sprintf("invalid minimum / maximum specified: maximum value of '%s' is below minimum value of '%s'", p.Maximum, p.Minimum),
				}
			}
		}
		return nil

	case NumericParameter:
		min := p.Minimum == "" || p.parsesAsNumber(p.Minimum)
		max := p.Maximum == "" || p.parsesAsNumber(p.Maximum)

		if !min && !max {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid minimum / maximum specified: '%s' / '%s' are not numbers", p.Minimum, p.Maximum),
			}
		}
		if !min {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid minimum specified: '%s' is not a number", p.Minimum),
			}
		}
		if !min {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid maximum specified: '%s' is not a number", p.Maximum),
			}
		}
		if p.Minimum != "" && p.Maximum != "" {
			floor, _ := p.parseAsNumber(p.Minimum)
			ceil, _ := p.parseAsNumber(p.Maximum)

			if ceil < floor {
				return MalformedParameterSpec{
					msg: fmt.Sprintf("invalid minimum / maximum specified: maximum value of '%s' is below minimum value of '%s'", p.Maximum, p.Minimum),
				}
			}
		}
		return nil

	case DurationParameter:
		min := p.Minimum == "" || p.parsesAsDuration(p.Minimum)
		max := p.Maximum == "" || p.parsesAsDuration(p.Maximum)

		if !min && !max {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid minimum / maximum specified: '%s' / '%s' are not durations", p.Minimum, p.Maximum),
			}
		}
		if !min {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid minimum specified: '%s' is not a duration", p.Minimum),
			}
		}
		if !min {
			return MalformedParameterSpec{
				msg: fmt.Sprintf("invalid maximum specified: '%s' is not a duration", p.Maximum),
			}
		}
		if p.Minimum != "" && p.Maximum != "" {
			floor, _ := p.parseAsDuration(p.Minimum)
			ceil, _ := p.parseAsDuration(p.Maximum)

			if ceil < floor {
				return MalformedParameterSpec{
					msg: fmt.Sprintf("invalid minimum / maximum specified: maximum value of '%s' is below minimum value of '%s'", p.Maximum, p.Minimum),
				}
			}
		}
		return nil

	default:
		return MalformedParameterSpec{
			msg: fmt.Sprintf("unrecognized parameter type '%s'", p.Type),
		}
	}
}

func (p Parameter) Validate(value string) error {
	if err := p.Wellformed(); err != nil {
		return err
	}

	switch p.Type {
	case TextParameter:
		return nil

	case BytesParameter:
		min, _ := p.parseAsBytes(p.Minimum)
		max, _ := p.parseAsBytes(p.Maximum)
		v, err := p.parseAsBytes(value)

		if err != nil {
			return err
		}
		if v < min {
			return ParameterValueOutOfRange{
				value: value,
				msg:   fmt.Sprintf("value '%s' is below minimum of '%s'", value, p.Minimum),
			}
		}
		if v > max {
			return ParameterValueOutOfRange{
				value: value,
				msg:   fmt.Sprintf("value '%s' is above maximum of '%s'", value, p.Maximum),
			}
		}
		return nil

	case NumericParameter:
		return nil

	case DurationParameter:
		return nil

	default:
		return nil
	}
}

func (p Parameter) parseAsBytes(in string) (int64, error) {
	for i, r := range in {
		if r >= '0' && r <= '9' || r == '.' {
			continue
		}

		qty, err := strconv.ParseInt(in[:i], 10, 64)
		if err != nil {
			return 0, err
		}

		switch strings.ToLower(in[i:]) {
		case "b":
			return qty, nil

		case "k", "kb", "ki", "kib":
			return qty * 1024, nil

		case "m", "mb", "mi", "mib":
			return qty * 1024 * 1024, nil

		case "g", "gb", "gi", "gib":
			return qty * 1024 * 1024 * 1024, nil

		case "t", "tb", "ti", "tib":
			return qty * 1024 * 1024 * 1024 * 1024, nil

		default:
			return 0, fmt.Errorf("'%s' is not a size", in)
		}
	}

	// unit-less size specification implies 'b'
	return strconv.ParseInt(in, 10, 64)
}

func (p Parameter) parsesAsBytes(in string) bool {
	_, err := p.parseAsBytes(in)
	return err == nil
}

func (p Parameter) parseAsNumber(in string) (int64, error) {
	return strconv.ParseInt(in, 10, 64)
}

func (p Parameter) parsesAsNumber(in string) bool {
	_, err := p.parseAsNumber(in)
	return err == nil
}

func (p Parameter) parseAsDuration(in string) (int64, error) {
	for i, r := range in {
		if r >= '0' && r <= '9' || r == '.' {
			continue
		}

		qty, err := strconv.ParseInt(in[:i], 10, 64)
		if err != nil {
			return 0, err
		}

		switch strings.ToLower(in[i:]) {
		case "s":
			return qty, nil

		case "m":
			return qty * 60, nil

		case "h":
			return qty * 60 * 60, nil

		case "d":
			return qty * 60 * 60 * 24, nil

		case "w":
			return qty * 60 * 60 * 24 * 7, nil

		case "y":
			return qty * 60 * 60 * 24 * 365, nil

		default:
			return 0, fmt.Errorf("'%s' is not a duration", in)
		}
	}

	// unit-less size specification implies 's'
	return strconv.ParseInt(in, 10, 64)
}

func (p Parameter) parsesAsDuration(in string) bool {
	_, err := p.parseAsDuration(in)
	return err == nil
}
