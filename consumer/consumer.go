package consumer

import (
	"fmt"

	"github.com/dius/libpact/api"
	"github.com/dius/libpact/mock"
	"github.com/dius/libpact/pactfile"
)

type Mux struct {
	Pact          *pactfile.Root
	CurrentServer *mock.Server
}

func (c *Mux) HandlerByName(name string) api.HandlerFunc {
	switch name {
	case "add_interaction", "POST /interaction":
		return c.AddInteraction
	case "mock", "POST /mock":
		return c.StartMockServer
	case "setup", "POST /setup":
		return c.Setup
	case "report", "GET /report":
		return c.Report
	case "save", "POST /save":
		return c.Save
	default:
		return api.NotFound
	}
}

// Setup resets the pactfile to empty ( / Creates a new pactfile )
func (c *Mux) Setup(req api.Request) (api.Response, error) {
	*c.Pact = *pactfile.New()
	return api.BuildStringResponse(200, "OK"), nil
}

// Report returns coverage reports after running all calls.
func (c *Mux) Report(req api.Request) (api.Response, error) {
	return nil, fmt.Errorf("Not Implemented")
}

// Save saves and uploads the pactfile
func (c *Mux) Save(req api.Request) (api.Response, error) {
	var config = &struct {
		Dest string `json:"dst"`
	}{}
	err := req.ReadBodyInto(config)
	if err != nil {
		return nil, err
	}
	err = pactfile.Save(config.Dest, c.Pact)
	if err != nil {
		return nil, err
	}
	return api.BuildStringResponse(200, "OK"), nil
}

func (c *Mux) StartMockServer(req api.Request) (api.Response, error) {

	if c.CurrentServer != nil {
		c.CurrentServer.Close()
	}

	var config = &struct {
		Bind string `json:"bind"`
	}{}
	err := req.ReadBodyInto(config)
	if err != nil {
		return nil, err
	}

	// Setup default parameters if not set
	if len(config.Bind) < 1 {
		//TODO: Dynamic High Ports
		config.Bind = "localhost:8080"
	}

	c.CurrentServer, err = mock.Serve(config.Bind)
	if err != nil {
		return nil, err
	}

	return api.BuildObjectResponse(200, config), nil
}

// CreateInteraction from the client adds a new interaction to the pactfile
func (c *Mux) AddInteraction(req api.Request) (api.Response, error) {
	interaction := pactfile.Interaction{}
	err := req.ReadBodyInto(&interaction)
	if err != nil {
		fmt.Printf("E1: %s\n", err.Error())
		return nil, err
	}

	c.CurrentServer.Interactions = append(c.CurrentServer.Interactions, interaction)
	c.Pact.Interactions = append(c.Pact.Interactions, interaction)

	return api.BuildStringResponse(200, "OK"), nil
}
