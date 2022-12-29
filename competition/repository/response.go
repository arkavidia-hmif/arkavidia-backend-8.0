package repository

type Response[T interface{}] struct {
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	URL     string `json:"URL,omitempty"`
}
