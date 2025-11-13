package persist

import (
	"encoding/json"
	"os"

	"matrixpulse/internal/engine"
)

type Persister struct {
	path string
	eng  *engine.Engine
}

func New(path string, eng *engine.Engine) *Persister {
	return &Persister{path: path, eng: eng}
}

func (p *Persister) Save() error {
	data := struct {
		Matrix interface{} `json:"matrix"`
		Mode   interface{} `json:"mode"`
		Alerts interface{} `json:"alerts"`
	}{
		Matrix: p.eng.Matrix(),
		Mode:   p.eng.Mode(),
		Alerts: p.eng.Alerts(),
	}

	f, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(data)
}
