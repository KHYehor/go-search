package model

type Position [2]int

type Progress struct {
	Finished   []string `json:"Finished"`
	Processing []string `json:"Processing"`
}
