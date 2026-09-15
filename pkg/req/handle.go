package req

import (
	"net/http"

	"erp-mes/pkg/res"
)

func HandleBody[T any](w http.ResponseWriter, req *http.Request) (*T, error) {
	body, err := Decode[T](req.Body)
	if err != nil {
		res.JSON(w, http.StatusBadRequest, "Couldn't parse JSON body")
		return body, err
	}

	err = IsValid(body)
	if err != nil {
		res.JSON(w, http.StatusBadRequest, err.Error())
		return body, err

	}

	return body, nil
}
