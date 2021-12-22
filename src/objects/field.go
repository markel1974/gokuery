package objects

type Nested struct {
	Path string `json:"path"`
}

type SubType struct {
	Nested *Nested `json:"nested"`
}

type Field struct {
	Name     string   `json:"name"`
	Value    string   `json:"value"`
	Lang     string   `json:"lang"`
	Type     string   `json:"type"`
	SubType  *SubType `json:"subType"`
	Scripted bool     `json:"scripted"`
	Script   string   `json:"script"`
}
