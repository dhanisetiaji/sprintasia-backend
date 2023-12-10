package utils

type ISearch interface {
	Get() *Find
	GetSearch() string
	GetStatus() string
}

type Find struct {
	Search string `query:"search"`
	Status string `query:"status"`
}

func (s *Find) Get() *Find {
	return s
}

func (s *Find) GetSearch() string {
	return s.Search
}

func (s *Find) GetStatus() string {
	return s.Status
}
