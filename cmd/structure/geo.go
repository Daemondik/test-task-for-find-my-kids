package structure

const (
	ReasonStatusOk = "OK"
	SourceFused    = "FUSED"
)

type Coordinate struct {
	Id        int     `json:"id"`
	ChildId   uint64  `json:"child_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"` // погрешность координаты в метрах (радиус)
	Accuracy  float64 `json:"accuracy"`
	Reason    string  `json:"reason"`
	Source    string  `json:"source"`
}
