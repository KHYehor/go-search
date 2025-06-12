package main

import (
	"github.com/gin-gonic/gin"
	"go-search/internal/adapter/http/handler"
	"go-search/internal/adapter/http/middleware"
	"go-search/internal/adapter/http/router"
	"go-search/internal/adapter/mongo"
	"go-search/internal/config"
	"go-search/internal/core/jobs"
	"go-search/internal/core/search"
	"go-search/internal/core/state"
	"net/http/pprof"
)

// Init profiler
func initProfiler(r *gin.Engine) {
	pprofGroup := r.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
		pprofGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
		pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
		pprofGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func main() {
	// Get config
	cfg := config.Load()

	// Db init
	db := mongo.NewDBManager(cfg)
	err := db.InitConnection()
	if err != nil {
		panic(err)
	}

	// Core business logic init
	f := search.NewIndexFactory()
	m := state.NewManager()
	j := jobs.NewJobs(cfg, m, db)

	// Init base
	r := gin.Default()

	//initProfiler(cfg)
	initProfiler(r)

	// Middlewares
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.UniqueRequestMiddleware())

	// Adapters (controllers) inits
	search := handler.NewSearch(f, j)
	jobs := handler.NewJobs(j)

	// Router
	router.NewRouter(r, search, jobs)

	// Listen the server
	r.Run(cfg.Addr)
}
