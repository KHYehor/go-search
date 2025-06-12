package router

import (
	"github.com/gin-gonic/gin"
	"go-search/internal/adapter/http/handler"
)

type Router struct {
	engine *gin.Engine
	s      *handler.Search
	j      *handler.Jobs
	prefix string
}

func NewRouter(e *gin.Engine, s *handler.Search, j *handler.Jobs) *Router {
	r := &Router{engine: e, s: s, j: j, prefix: "api"}
	r.setProcessRoutes()
	r.setJobsRoutes()

	return r
}

func (r *Router) setProcessRoutes() {
	domain := "/search"

	// Post /api/search/process
	r.engine.POST(r.prefix+domain+"/process", r.s.ProcessInput)
}

func (r *Router) setJobsRoutes() {
	domain := "/jobs"

	// GET /api/jobs/:id
	r.engine.GET(r.prefix+domain+"/:id", r.j.GetAJob)

	// GET /api/jobs
	r.engine.GET(r.prefix+domain, r.j.ListAllJobs)
}
