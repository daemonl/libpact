package api

import (
	"fmt"
	"net/http"

	"github.com/dius/libpact/pactfile"
)

type Consumer struct {
	Pact *pactfile.Root
}

func ConsumerServe(bind string, pact *pactfile.Root) error {
	api := &Consumer{
		Pact: pact,
	}

	mux := http.NewServeMux()
	mux.Handle("/interaction", Handler{POST: api.CreateInteraction})

	return http.ListenAndServe(bind, mux)
}

func (api *Consumer) HandleCall(call string, data string) FFIResponse {
	switch call {
	case "add_interaction":
		return HandleFFI(data, api.CreateInteraction)
	default:
		fmt.Println(call)
		return FFIResponse{
			Status:  404,
			Message: "Call not found",
		}
	}
}

// CreateInteraction from the client adds a new interaction to the pactfile
func (api *Consumer) CreateInteraction(req Request) (Response, error) {
	interaction := pactfile.Interaction{}
	err := req.ReadBodyInto(&interaction)
	if err != nil {
		fmt.Printf("E1: %s\n", err.Error())
		return nil, err
	}

	//TODO: Validate
	api.Pact.Interactions = append(api.Pact.Interactions, interaction)

	return GetStringResponse(200, "OK"), nil
}
