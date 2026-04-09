package db_manager_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arbhalerao/meerkat/utils"
)

const (
	MaxRetries     = 10
	MaxBackoff     = 30 * time.Second
	InitialBackoff = 1 * time.Second
)

type DBManagerClient struct {
	managerAddr string
}

func NewDBManagerClient(managerAddr, region string) *DBManagerClient {
	return &DBManagerClient{
		managerAddr: managerAddr,
	}
}

func (c *DBManagerClient) RegisterWithManager(region, grpcAddr string, ready chan<- bool) {
	backoff := InitialBackoff

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		data := map[string]string{
			"region":    region,
			"grpc_addr": grpcAddr,
		}
		payload, _ := json.Marshal(data)

		resp, err := http.Post(fmt.Sprintf("http://%s/register", c.managerAddr), "application/json", bytes.NewBuffer(payload))
		if err == nil && resp.StatusCode == http.StatusOK {
			utils.Logger.Info().Msgf("Successfully registered with db_manager (%s)", c.managerAddr)
			resp.Body.Close()
			ready <- true
			return
		}

		if resp != nil {
			resp.Body.Close()
		}

		utils.Logger.Warn().Msgf("Failed to register with db_manager (%s), attempt %d/%d, retrying in %v...",
			c.managerAddr, attempt, MaxRetries, backoff)

		time.Sleep(backoff)

		backoff *= 2
		if backoff > MaxBackoff {
			backoff = MaxBackoff
		}
	}

	utils.Logger.Error().Msgf("Failed to register with db_manager after %d attempts", MaxRetries)
	ready <- false
}
