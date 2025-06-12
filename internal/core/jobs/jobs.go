package jobs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-search/internal/config"
	"go-search/internal/entity"
	"go-search/internal/model"
	"os"
	"sync"
)

type Jobs struct {
	manager StateManager
	db      Db
	output  string
}

func NewJobs(cfg *config.Config, manager StateManager, db Db) *Jobs {
	return &Jobs{manager: manager, output: cfg.OutputDir, db: db}
}

func (j *Jobs) getActiveJobs() []string {
	return j.manager.GetAll()
}

func (j *Jobs) getFinishedJobs(c *gin.Context) ([]string, error) {
	jobs, err := j.db.GetAllItems(c)

	return jobs, err
}

func (j *Jobs) GetAllJobs(c *gin.Context) (*model.Progress, error) {
	res := model.Progress{
		Processing: j.getActiveJobs(),
	}
	f, err := j.getFinishedJobs(c)
	if err != nil {
		fmt.Println("Error getting finished jobs:", err)
		return nil, err
	}
	res.Finished = f

	return &res, nil
}

func (j *Jobs) CreateJob(key string) {
	j.manager.Add(key)
}

func (j *Jobs) saveToFile(wg *sync.WaitGroup, ch chan error, id string, idx Indexer) {
	defer wg.Done()

	err := os.MkdirAll(j.output, 0755)
	if err != nil {
		ch <- fmt.Errorf("failed to create directory: %w", err)
		return
	}

	// Construct the full file path
	path := j.output + id + ".json"

	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			ch <- fmt.Errorf("file already exists: %s", path)
			return
		}
		ch <- err
		return
	}
	defer file.Close()

	// Extract json
	json, err := idx.GetJson()
	if err != nil {
		ch <- fmt.Errorf("failed to marshal json: %w", err)
		return
	}
	// Write the JSON content
	_, err = file.Write(json)
	if err != nil {
		ch <- err
	}
}

func (j *Jobs) saveToDb(c *gin.Context, wg *sync.WaitGroup, ch chan error, id string, duration int64, idx Indexer) {
	defer wg.Done()

	err := j.db.InsertItem(c, id, duration, idx.GetResult())
	if err != nil {
		ch <- err
	}
}

func (j *Jobs) SaveJobResult(c *gin.Context, id string, duration int64, idx Indexer) (result []error) {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)

	// Concurrent execution
	go j.saveToFile(&wg, errCh, id, idx)
	go j.saveToDb(c, &wg, errCh, id, duration, idx)

	wg.Wait()
	close(errCh)
	for e := range errCh {
		result = append(result, e)
	}

	// Remove job from manager
	j.manager.Remove(id)

	return result
}

func (j *Jobs) GetAJob(c *gin.Context, id string) (*entity.SearchResult, error) {
	job, err := j.db.GetItem(c, id)
	if err != nil {
		fmt.Printf("Error getting item %s: %v", id, err)
		return nil, err
	}

	return job, nil
}
