package provider

import (
	"github.com/dius/libpact/api"
	"github.com/dius/libpact/pactfile"
)

type Mux struct {
	Pact         *pactfile.Root
	Interactions []Interaction
}

type Interaction struct {
	pactfile.Interaction
}

func (m *Mux) HandlerByName(name string) api.HandlerFunc {
	switch name {
	case "start", "POST /start":
		return m.Start
	case "run", "POST /run":
		return c.Run
	default:
		return api.NotFound
	}
}

type StepResponse struct {
	pactfile.DiffReport `json:"diff,omitempty"`
	Next                pactfile.ProviderInteractionSpec `json:"next"`
}

func (resp *StepResponse) StatusCode() int {
	return 200
}
func (resp *StepResponse) GetEncodable() interface{} {
	return resp
}

func (m *Mux) Start(req api.Request) (api.Response, error) {
	var config = &struct {
		Pactfile string `json:"pactfile"`
	}{}
	err := req.ReadBodyInto(config)
	if err != nil {
		return nil, err
	}

	// Load Pactfile
	err, pact := pactfile.LoadFile(config.Pactfile)
	if err != nil {
		return nil, err
	}
	m.Pact = pact
	m.Interactions = make([]Interaction, len(m.Pact.Interactions), len(m.Pact.Interactions))
	for i, interaction := range m.Pact.Interactions {
		m.Interactions[i] = Interaction{
			Interaction: interaction,
		}
	}

	resp := StepResponse{}
	resp.Next = m.Pact.NextAfter("")

	return resp, nil
}

func (m *Mux) Run(req api.Request) (api.Response, error) {
	var spec = &struct {
		InteractionID int    `json:"interactionId"`
		HTTPRoot      string `json:"httpRoot"`
	}{}
	err := req.ReadBodyInto(spec)
	if err != nil {
		return nil, err
	}

	resp := &StepResponse{}

	interaction := m.Pact.GetInteraction(spec.InteractionID)
	resp.Diff = interaction.Run(spec.HTTPRoot)

	resp.Next = m.Pact.NextAfter(spec.InteractionID)
	return resp, nil
}
