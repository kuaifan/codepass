package app

type MultipassModel struct {
	Ip   string
	Port string
}

type infoModel struct {
	Errors any            `json:"errors"`
	Info   map[string]any `json:"info"`
}

var (
	MConf MultipassModel
)
