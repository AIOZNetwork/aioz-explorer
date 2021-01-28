package domain

import "net/http"

type LCDRepository interface {

}

type LCDUsecase interface {
	Request(host, port, method, path string, payload []byte) (*http.Response, string, error)
	ExtractResultFromResponse(body []byte) ([]byte, error)
}
