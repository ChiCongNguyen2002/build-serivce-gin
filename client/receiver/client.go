package receiver

import (
	"build-service-gin/common/client"
	"build-service-gin/common/utils"
	"build-service-gin/config"
	"build-service-gin/internal/domains"
	"build-service-gin/pkg/helpers/constants"
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	pathPostOrder = "/rewards/receiver/v1/internal/order"
)

type ReceiverClient struct {
	conf *config.SystemConfig
	cl   *client.Client
}

type IReceiverClient interface {
	PostOrder(ctx context.Context, message domains.OrderMessage) (*OrderResp, error)
}

var (
	instanceReceiverClient *ReceiverClient
	onceReceiverClient     sync.Once
)

func NewReceiverClient() IReceiverClient {
	onceReceiverClient.Do(func() {
		conf := config.GetInstance()
		cl := client.NewClient(conf.RewardIntegrationUrl, 30*time.Second, 0, 1*time.Second, nil)
		instanceReceiverClient = &ReceiverClient{
			conf: conf,
			cl:   cl,
		}
	})

	return instanceReceiverClient
}

func (c *ReceiverClient) PostOrder(ctx context.Context, message domains.OrderMessage) (*OrderResp, error) {
	region := ctx.Value(utils.KeyRegion).(string)
	headers := map[string]string{
		"Authorization":   fmt.Sprintf("Bearer %s", c.conf.InternalToken),
		"X-Client-Region": region,
	}
	// Create action and source send order
	action := constants.ActionEarn
	source := constants.SourceRewards
	apiURL := fmt.Sprintf("%s/%s/%s", pathPostOrder, action, source)

	var res OrderResp
	_, err := c.cl.R().SetContext(ctx).SetHeaders(headers).SetBody(message).SetResult(&res).Post(apiURL)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
