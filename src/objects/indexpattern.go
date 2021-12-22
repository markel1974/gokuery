package objects

import "errors"

type IndexPattern struct {
	fields []*Field
}

func NewIndexPattern() *IndexPattern{
	ip := &IndexPattern{
		fields: nil,
	}
	return ip
}

func (ip * IndexPattern) AddField(field * Field) {
	if field == nil {
		return
	}
	ip.fields = append(ip.fields, field)
}

func (ip * IndexPattern) FieldsLen() int {
	return len(ip.fields)
}

func (ip * IndexPattern) Find(fieldName string) *Field{
	for _, f := range ip.fields {
		if f.Name == fieldName {
			return f
		}
	}
	return nil
}

func (ip * IndexPattern) Filter(fn func(string) bool) []*Field{
	var out []*Field
	for _, f := range ip.fields {
		if fn(f.Name) {
			out = append(out, f)
		}
	}
	return out
}

func (ip * IndexPattern) VerifyFields(nestedPath string) []error {
	var errs []error
	for _, field := range ip.fields {
		var nestedPathFromField string
		if field.SubType != nil && field.SubType.Nested != nil && len(field.SubType.Nested.Path) > 0 {
			nestedPathFromField = field.SubType.Nested.Path
		}
		if len(nestedPath) > 0 && len(nestedPathFromField) == 0 {
			t := field.Name + " is not a nested field but is in nested group" + nestedPath + "in the KQL expression"
			errs = append(errs, errors.New(t))
			continue
		}
		if len(nestedPathFromField) > 0 && len(nestedPath) == 0 {
			t := field.Name + " is a nested field, but is not in a nested group in the KQL expression"
			errs = append(errs, errors.New(t))
			continue
		}
		if nestedPathFromField != nestedPath {
			t := field.Name +  " is being queried with the incorrect nested path. The correct path is " + nestedPathFromField
			errs = append(errs, errors.New(t))
			continue
		}
	}
	return errs
}