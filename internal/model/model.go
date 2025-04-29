package model

type WebResponse[T any] struct {
	Data   T      `json:"data"`
	Errors string `json:"errors"`
}

type WebResponses[T any] struct {
	Data   *[]T   `json:"data"`
	Errors string `json:"errors"`
}
