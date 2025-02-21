package tf

import (
	"errors"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type variables struct {
	v map[string]any
	m map[string]hcl.Expression
}

func newVariables() variables {
	return variables{
		v: map[string]any{
			"string": cty.StringVal("string"),
		},
		m: map[string]hcl.Expression{},
	}
}

func (v variables) setExpr(keys []string, e hcl.Expression) {
	key := strings.Join(keys, ".")
	v.m[key] = e
}

func (v variables) set(keys []string, vv cty.Value) error {
	if len(keys) == 0 {
		return errors.New("keys is empty")
	}

	current := v.v
	for i, key := range keys {
		if i == len(keys)-1 {
			current[key] = vv
			return nil
		}

		val, ok := current[key]
		if !ok {
			current[key] = map[string]any{}
			val = current[key]
		}

		current, ok = val.(map[string]any)
		if !ok {
			return errors.New("unexpected error")
		}
	}

	return errors.New("unexpected error")
}

func (v variables) ctx() (*hcl.EvalContext, error) {
	i := 0
	try := 10
	var ctx *hcl.EvalContext
	for {
		vars := toCtyMap(v.v).AsValueMap()
		ctx = &hcl.EvalContext{
			Variables: vars,
		}
		var err error
		for k, e := range v.m {
			err = nil
			vv, diags := e.Value(ctx)
			if diags.HasErrors() {
				err = diags
				if i == try {
					return nil, diags
				}
			}
			// set value
			keys := strings.Split(k, ".")
			if err := v.set(keys, vv); err != nil {
				return nil, err
			}
		}
		if err == nil {
			break
		}
		i++
	}
	return ctx, nil
}

func toCtyMap(m map[string]any) cty.Value {
	obj := map[string]cty.Value{}
	for k, v := range m {
		switch vv := v.(type) {
		case map[string]any:
			obj[k] = toCtyMap(vv)
		case cty.Value:
			obj[k] = vv
		}
	}
	return cty.ObjectVal(obj)
}
