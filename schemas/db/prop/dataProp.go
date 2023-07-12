package schemas

type DataProp struct {
	Length      uint `json:"length"`
	TotalPage   uint `json:"total_page"`
	CurrentPage uint `json:"current_page"`
}
