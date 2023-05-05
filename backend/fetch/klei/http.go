package klei

// WrappedResponse is Klei's special response format where the actual
// data is wrapped in a "GET" json field.
type WrappedResponse[T any] struct {
	GET T `json:"GET"`
}
