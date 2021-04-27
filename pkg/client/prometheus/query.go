package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/common/model"
)

type Query interface {
	IsActivity(ctx context.Context) (bool, error)
}

func (pc Client) IsActivity(ctx context.Context) (bool, error) {
	value, _, err := pc.Query(ctx, "sum(rate(veidemann_page_requests_total[5m]))", time.Now())
	if err != nil {
		return false, err
	}
	switch value.Type() {
	case model.ValVector:
		if vector, ok := value.(model.Vector); !ok {
			return false, nil
		} else if len(vector) == 0 {
			return false, fmt.Errorf("expected vector to have values: %#v", vector)
		} else {
			return vector[0].Value > 0, nil
		}
	default:
		return false, nil
	}
}
