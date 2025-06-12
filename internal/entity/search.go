package entity

type SearchResult struct {
	Id         string `bson:"_id"`
	DurationMs int64  `bson:"durationMs"`
	Data       any    `bson:"data"`
}
