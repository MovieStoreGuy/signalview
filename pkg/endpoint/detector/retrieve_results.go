package detector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	internalClient "github.com/MovieStoreGuy/signalview/pkg/client"
	"github.com/MovieStoreGuy/signalview/pkg/config"
	"github.com/sirupsen/logrus"
)

const (
	domainFormat = `https://api.%s.signalfx.com/v2/detector`
)

// QueryParameter defines what can be passed in as value to fetch GetMatching detectors
type QueryParameter struct {
	Limit  int
	Name   string
	Offset int
	Tags   string
}

func (qp *QueryParameter) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("?limit=%d", qp.Limit))
	sb.WriteString(fmt.Sprintf("&offset=%d", qp.Offset))
	if qp.Name != "" {
		sb.WriteString(fmt.Sprintf("&name=%s", qp.Name))
	}
	if qp.Tags != "" {
		sb.WriteString(fmt.Sprintf("&tags=%s", qp.Tags))
	}
	return sb.String()
}

// GetMatching returns all v2 detectors that match the provided query, if the query is not defined, it will return all detectors.
func GetMatching(ctx context.Context, log *logrus.Logger, client *http.Client, conf *config.Runtime, query *QueryParameter) (*BundledPayload, error) {
	if client == nil {
		return nil, errors.New("undefined client passed")
	}
	if conf == nil {
		return nil, errors.New("undefined config passed")
	}
	if query == nil {
		query = &QueryParameter{
			Limit:  80,
			Offset: 0,
		}
	}
	domain := fmt.Sprintf(domainFormat, conf.Relm)
	requestor := internalClient.NewCachedRequest(conf.Token)

	// Loop through all detectors that match the query and go from there
	completeBundle := &BundledPayload{}
	for {
		req, err := requestor(ctx, http.MethodGet, domain+query.String(), nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(resp.Status)
		}
		buff, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var bundle BundledPayload
		if err = json.Unmarshal(buff, &bundle); err != nil {
			log.Info(string(buff))
			return nil, err
		}
		if completeBundle.Count == 0 {
			completeBundle.Count = bundle.Count
		}
		completeBundle.Results = append(completeBundle.Results, bundle.Results...)
		if completeBundle.Count == int64(len(completeBundle.Results)) {
			break
		}
		query.Offset += query.Limit
	}
	return completeBundle, nil
}
