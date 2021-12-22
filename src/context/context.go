package context

type Nested struct {
	Path string `json:"path"`
}

func NewNested() * Nested {
	return &Nested{}
}

func (n * Nested) Clone() *Nested {
	z := NewNested()
	z.Path = n.Path
	return z
}

type Context struct {
	Nested * Nested `json:"nested"`
}

func NewContext() * Context {
	return &Context{}
}

func (ctx * Context) Clone() *Context {
	z := NewContext()
	if ctx.Nested != nil {
		z.Nested = ctx.Nested.Clone()
	}
	return z
}