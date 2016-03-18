package provider

// NextInteraction from the provider should be called in a loop from the
// provider until there are none left. It returns the metadata

// RunInteraction from the provider triggers a request back to the provider.
// When the provider responds, the response is compared with the matcher, and
// this request returns the result
