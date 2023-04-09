package db

import (
	"context"
	"fmt"
	"time"

	influxdbV2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Client struct {
	Client   influxdbV2.Client
	WriteAPI api.WriteAPIBlocking
	QueryAPI api.QueryAPI
}

func NewClient(baseURL, token, org, bucket string) (*Client, error) {
	client := influxdbV2.NewClient(baseURL, token)
	writeAPI := client.WriteAPIBlocking(org, bucket)
	queryAPI := client.QueryAPI(org)

	return &Client{
		Client:   client,
		WriteAPI: writeAPI,
		QueryAPI: queryAPI,
	}, nil
}

func (c *Client) WriteData(data interface{}, tags map[string]string, fields map[string]interface{}, time time.Time) error {
	point := influxdbV2.NewPoint(fmt.Sprintf("%T", data), tags, fields, time)
	err := c.WriteAPI.WritePoint(context.Background(), point)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryData(query string) (*api.QueryTableResult, error) {
	return c.QueryAPI.Query(context.Background(), query)
}

func (c *Client) Close() {
	c.Client.Close()
}
