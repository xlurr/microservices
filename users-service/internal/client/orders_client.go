package client

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type OrdersServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// создаём клиента берем урл из енв
func NewOrdersServiceClient() *OrdersServiceClient {
	ordersURL := os.Getenv("ORDERS_SERVICE_URL")
	if ordersURL == "" {
		ordersURL = "http://orders-service:8082"
	}

	return &OrdersServiceClient{
		baseURL: ordersURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// DeleteUserOrders удаляет все заказы пользователя по иду
func (c *OrdersServiceClient) DeleteUserOrders(userID int64) error {
	url := fmt.Sprintf("%s/api/orders/user/%d", c.baseURL, userID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete user orders: %w", err)
	}
	defer resp.Body.Close()

	// 204 No Content или 200 OK — оба успешны
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete user orders: status %d", resp.StatusCode)
	}

	return nil
}
