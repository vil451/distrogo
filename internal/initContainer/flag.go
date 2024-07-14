package initContainer

type Flags struct {
	Name *string
}

func NewFlags() *Flags {
	return &Flags{
		Name: strPtr("name"),
	}
}

func strPtr(s string) *string {
	return &s
}
