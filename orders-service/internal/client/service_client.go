package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type ServiceClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewServiceClient(baseURL string) *ServiceClient {
    return &ServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *ServiceClient) CreatePayment(orderID int64, amount float64) error {
    payload := map[string]interface{}{
        "order_id": orderID,
        "amount":   amount,
        "status":   "pending",
    }
    
    body, _ := json.Marshal(payload)
    
    req, err := http.NewRequest("POST", 
        fmt.Sprintf("%s/api/payments", c.baseURL), 
        bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("failed to create payment: status %d", resp.StatusCode)
    }
    
    return nil
}
