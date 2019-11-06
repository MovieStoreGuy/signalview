package alerts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MovieStoreGuy/signalview/pkg/client"
	"github.com/MovieStoreGuy/signalview/pkg/config"
	"github.com/sirupsen/logrus"
)

const (
	mutingDomain = `https://api.%s.signalfx.com/v2/alertmuting`
)

// QueryMutingRules defines what can be passed in as value to fetch the matching muting rules
type QueryMutingRules struct {
	Include string
	Limit   int32
	Offset  int32
	OrderBy string
	Query   string
}

func (qp *QueryMutingRules) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("?limit=%d", qp.Limit))
	sb.WriteString(fmt.Sprintf("&offset=%d", qp.Offset))
	if qp.Include != "" {
		sb.WriteString(fmt.Sprintf("&include=%s", qp.Include))
	}
	if qp.OrderBy != "" {
		sb.WriteString(fmt.Sprintf("&orderby=%s", qp.OrderBy))
	}
	if qp.Query != "" {
		sb.WriteString(fmt.Sprintf("&quert=%s", qp.Query))
	}
	return sb.String()
}

func GetMutingRules(ctx context.Context, log *logrus.Logger, c *http.Client, conf *config.Runtime, query *QueryMutingRules) (*BundledMutingRuleResults, error) {
	if c == nil {
		return nil, errors.New("http client not defined")
	}
	if conf == nil {
		return nil, errors.New("configu was not defined")
	}
	if query == nil {
		query = &QueryMutingRules{
			Limit:  80,
			Offset: 0,
		}
	}
	domain := fmt.Sprintf(mutingDomain, conf.Relm)
	requestor := client.NewCachedRequest(conf.Token)

	completeBundle := &BundledMutingRuleResults{}
	for {
		req, err := requestor(ctx, http.MethodGet, domain+query.String(), nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.Do(req)
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
		var bundle BundledMutingRuleResults
		if err = json.Unmarshal(buff, &bundle); err != nil {
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
