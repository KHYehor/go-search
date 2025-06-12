package jobs

import (
	"github.com/gin-gonic/gin"
	"go-search/internal/entity"
	"go-search/internal/model"
)

type StateManager interface {
	Add(key string)
	Remove(key string)
	GetAll() []string
}

type Db interface {
	InsertItem(c *gin.Context, id string, duration int64, item any) error
	GetItem(c *gin.Context, id string) (*entity.SearchResult, error)
	GetAllItems(c *gin.Context) ([]string, error)
}

type Indexer interface {
	Close()
	Search(words []string)
	GetJson() ([]byte, error)
	GetResult() map[string][]model.Position
}
