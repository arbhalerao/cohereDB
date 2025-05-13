package db_manager_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arbhalerao/cohereDB/utils"
)

type DBManagerClient struct {
	managerAddr string
}

func NewDBManagerClient(managerAddr, region string) *DBManagerClient {
	return &DBManagerClient{
		managerAddr: managerAddr,
	}
}

// RegisterWithManager attempts to register the DB server with db_manager and signals success.
func (c *DBManagerClient) RegisterWithManager(region, grpcAddr string, ready chan<- bool) {
	retries := 0
	for {
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

		utils.Logger.Warn().Msgf("Failed to register with db_manager (%s), retrying... (%d)", c.managerAddr, retries+1)
		retries++
		time.Sleep(time.Duration(retries) * time.Second) // Exponential backoff
	}
}
