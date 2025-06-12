package handler

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"go-search/internal/core/jobs"
	"go-search/internal/entity"
	"go-search/internal/model"
)

type IndexerFactory interface {
	CreateNewIndex(scanner *bufio.Scanner) jobs.Indexer
}

type JobsManager interface {
	CreateJob(key string)
	GetAllJobs(c *gin.Context) (*model.Progress, error)
	SaveJobResult(c *gin.Context, id string, duration int64, idx jobs.Indexer) []error
	GetAJob(c *gin.Context, id string) (*entity.SearchResult, error)
}
