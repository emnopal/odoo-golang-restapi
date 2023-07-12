package schemas

type GetResponse struct {
	Status      int         `json:"status"`
	Message     string      `json:"message"`
	Length      uint        `json:"length"`
	TotalPage   uint        `json:"total_page"`
	CurrentPage uint        `json:"current_page"`
	Data        interface{} `json:"data"`
}

type PostResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
