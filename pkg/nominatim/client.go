package nominatim

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	server     string
	user       string
	password   string
	httpClient *http.Client
}

func NewClient(server, user, password string) *Client {
	return &Client{
		server:     strings.TrimRight(server, "/"),
		user:       user,
		password:   password,
		httpClient: &http.Client{},
	}
}

type reverseResponse struct {
	DisplayName string `json:"display_name"`
	Error       string `json:"error"`
}

type searchResult struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

// ReverseGeocode translates coordinates into a human-readable address.
func (c *Client) ReverseGeocode(lat, lon float64) (string, error) {
	query := url.Values{}
	query.Set("format", "jsonv2")
	query.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	query.Set("lon", strconv.FormatFloat(lon, 'f', -1, 64))

	var result reverseResponse
	if err := c.get("/reverse", query, &result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", fmt.Errorf("nominatim: %s", result.Error)
	}

	return result.DisplayName, nil
}

// Geocode translates a free-text address into coordinates. When Nominatim
// returns multiple matches, the first (its best-ranked) result is used.
func (c *Client) Geocode(address string) (lat, lon float64, err error) {
	query := url.Values{}
	query.Set("format", "jsonv2")
	query.Set("q", address)

	var results []searchResult
	if err := c.get("/search", query, &results); err != nil {
		return 0, 0, err
	}
	if len(results) == 0 {
		return 0, 0, fmt.Errorf("nominatim: no results for address %q", address)
	}

	lat, err = strconv.ParseFloat(results[0].Lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("nominatim: invalid latitude in response: %w", err)
	}
	lon, err = strconv.ParseFloat(results[0].Lon, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("nominatim: invalid longitude in response: %w", err)
	}

	return lat, lon, nil
}

func (c *Client) get(path string, query url.Values, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.server+path+"?"+query.Encode(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.user, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nominatim: unexpected status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
