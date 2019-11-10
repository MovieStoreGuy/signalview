package detector_test

import (
	"fmt"
	"testing"

	"github.com/MovieStoreGuy/signalview/pkg/endpoint/detector"
	"github.com/stretchr/testify/require"
)

func TestQueryBuilder(t *testing.T) {
	const (
		defaultQueryResultFormat = `?limit=%d&offset=%d`
		nameParamAdded           = `&name=%s`
		tagsParamAdded           = `&tags=%s`
	)
	t.Parallel()
	qp := &detector.QueryParameter{
		Limit:  10,
		Offset: 7,
	}
	require.Equal(t, fmt.Sprintf(defaultQueryResultFormat, qp.Limit, qp.Offset), fmt.Sprint(qp))
	qp.Name = "fart"
	require.Equal(t, fmt.Sprintf(defaultQueryResultFormat+nameParamAdded, qp.Limit, qp.Offset, qp.Name), fmt.Sprint(qp))
	qp.Name = ""
	require.Equal(t, fmt.Sprintf(defaultQueryResultFormat, qp.Limit, qp.Offset), fmt.Sprint(qp))
	qp.Tags = "derp"
	require.Equal(t, fmt.Sprintf(defaultQueryResultFormat+tagsParamAdded, qp.Limit, qp.Offset, qp.Tags), fmt.Sprint(qp))
}
