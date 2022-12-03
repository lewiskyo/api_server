package entity

type ParseRq struct {
	Name string `form:"name"`
	Age  uint32 `form:"age"`
}

type ParseRb struct {
	Data  []string `json:"data"`
	Data2 []struct {
		Image string `json:"image"`
		Size  uint32 `json:"size"`
	}
}

type ParseRet struct {
	Status uint        `json:"status"`
	Data   interface{} `json:"data"`
}
