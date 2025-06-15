package model

type Position [2]uint32

type Progress struct {
	Finished   []string `json:"Finished"`
	Processing []string `json:"Processing"`
}
