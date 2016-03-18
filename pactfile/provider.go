package pactfile

// ProviderInteractionSpec is returned to the provider to tell it how to
// prepare for and call the 'next' interaction
type ProviderInteractionSpec struct {
	ProviderState string `json:"provider_state"`
	ID            string `json:"id"`
}
