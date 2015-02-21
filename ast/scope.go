package ast

type Scope struct {
	Outer   *Scope
	Objects map[string]*Object
}

func NewScope(outer *Scope) *Scope {
	return &Scope{Outer: outer, Objects: make(map[string]*Object)}
}

func (s *Scope) Lookup(name string) *Object {
	scope := s
	for scope != nil {
		if obj, ok := scope.Objects[name]; ok {
			return obj
		}
		scope = scope.Outer
	}
	return nil
}

func (s *Scope) Insert(name string, obj *Object) {
	s.Objects[name] = obj
}
