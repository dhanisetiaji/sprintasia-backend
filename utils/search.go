package utils

type ISearch interface {
	Get() *Find
	GetSearch() string
}

type Find struct {
	Search string `query:"search"`
}

func (s *Find) Get() *Find {
	return s
}

func (s *Find) GetSearch() string {
	return s.Search
}
