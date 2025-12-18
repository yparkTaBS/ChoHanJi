package CreateRoom

type Request struct {
	MapWidth  int    `json:"MapWidth" validate:"required,gt=0"`
	MapHeight int    `json:"MapHeight" validate:"required,gt=0"`
	Items     string `json:"Items" validate:"required"`
}
