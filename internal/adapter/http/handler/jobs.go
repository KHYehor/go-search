package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Jobs struct {
	jobs JobsManager
}

func NewJobs(j JobsManager) *Jobs {
	return &Jobs{jobs: j}
}

func (j *Jobs) ListAllJobs(c *gin.Context) {
	res, err := j.jobs.GetAllJobs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (j *Jobs) GetAJob(c *gin.Context) {
	id := c.Param("id")

	res, err := j.jobs.GetAJob(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
