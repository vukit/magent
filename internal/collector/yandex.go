package collector

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/vukit/magent/internal/metric"

	"github.com/golang-jwt/jwt/v4"
)

type Metric struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type Metrics struct {
	TS     string `json:"ts"`
	Labels struct {
		Host string `json:"host"`
	} `json:"labels"`
	Metrics []Metric `json:"metrics"`
}

type Yandex Collector

func (c *Yandex) Send(metrics metric.Metrics, parameters map[string]interface{}) (err error) {
	query := url.Values{
		"folderId": {parameters["folderId"].(string)},
		"service":  {"custom"},
	}

	reqURL := "https://monitoring.api.cloud.yandex.net/monitoring/v2/data/write/?" + query.Encode()

	token, err := c.getIAMToken(parameters["iss"].(string), parameters["kid"].(string))
	if err != nil {
		return err
	}

	body, err := c.getBodyData(metrics)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", reqURL, body)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return err
	}

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		resp.Body.Close()

		return err
	}

	resp.Body.Close()

	return nil
}

func (c *Yandex) getBodyData(metrics metric.Metrics) (*bytes.Buffer, error) {
	data := Metrics{}
	data.TS = time.Now().Format("2006-01-02T15:04:05Z")
	data.Labels.Host = c.Config.HostName

	for _, metric := range metrics {
		switch metric.TValue {
		case "IGAUGE":
			data.Metrics = append(data.Metrics, Metric{Name: metric.Name, Type: metric.TValue, Value: metric.IValue})
		case "DGAUGE":
			data.Metrics = append(data.Metrics, Metric{Name: metric.Name, Type: metric.TValue, Value: metric.DValue})
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	bodyData := new(bytes.Buffer)

	_, err = bodyData.Write(jsonData)
	if err != nil {
		return nil, err
	}

	return bodyData, nil
}

func (c *Yandex) getIAMToken(iss, kid string) (string, error) {
	token, err := c.signedToken(iss, kid)
	if err != nil {
		return "", err
	}

	body := strings.NewReader(fmt.Sprintf(`{"jwt":%q}`, token))

	req, err := http.NewRequestWithContext(context.Background(), "POST", "https://iam.api.cloud.yandex.net/iam/v1/tokens", body)
	if err != nil {
		return "", nil
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		resp.Body.Close()

		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return "", fmt.Errorf("%s: %s", resp.Status, body)
	}

	var data struct {
		IAMToken string `json:"iamToken"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.IAMToken, nil
}

func (c *Yandex) signedToken(iss, kid string) (signed string, err error) {
	claims := jwt.RegisteredClaims{
		Issuer:    iss,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  []string{"https://iam.api.cloud.yandex.net/iam/v1/tokens"},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = kid

	privateKey, err := c.loadPrivateKey()
	if err != nil {
		return "", err
	}

	signed, err = token.SignedString(privateKey)
	if err != nil {
		return "", nil
	}

	return signed, nil
}

func (c *Yandex) loadPrivateKey() (rsaPrivateKey *rsa.PrivateKey, err error) {
	data, err := os.ReadFile(c.Config.PrivateKey)
	if err != nil {
		return nil, err
	}

	rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, err
	}

	return rsaPrivateKey, nil
}
