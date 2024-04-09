package repository

type Repository interface {
	GetPages() ([]string, error)
	GetOutboundLinks(pageId string) ([]string, error)
	SetPageRank(pageId string, rank float64) error
}
