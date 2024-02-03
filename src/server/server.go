package server

import (
	"access-tester/src/scraper"
	"context"
	"fmt"
	url2 "net/url"
)

type ServerImpl struct {
	Scraper *scraper.Scraper
}

func (s ServerImpl) CheckStatus(ctx context.Context, request CheckStatusRequestObject) (CheckStatusResponseObject, error) {
	url, err := url2.Parse(fmt.Sprintf("http://%s", request.Params.Host))
	if err != nil {
		return CheckStatus400Response{}, nil
	}
	result, err := s.Scraper.CheckIfBlocked(url.String())
	if err != nil {
		return CheckStatus500Response{}, nil
	}
	result.Host = request.Params.Host
	return CheckStatus200JSONResponse(result), nil
}

var _ StrictServerInterface = (*ServerImpl)(nil)
