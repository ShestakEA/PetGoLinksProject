package request

import (
	"go/adv-demo/pkg/response"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, req *http.Request) (*T, error) {
	body, err := Decode[T](req.Body)
	if err != nil {
		response.Response(*w, err.Error(), 402)
		return nil, err
	}
	/*err = IsValid(body)
	if err != nil {
		response.Response(*w, err.Error(), 402)
		return nil, err
	}*/
	return &body, nil
}
