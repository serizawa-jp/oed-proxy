package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const (
	baseURL     = "https://od-api.oxforddictionaries.com/api/v2"
	defaultLang = "en-us"
)

func mustGetBaseURL() *url.URL {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	return u
}

type OxfordClient struct {
	httpClient *http.Client
}

func NewOxfordClient(httpClient *http.Client) *OxfordClient {
	return &OxfordClient{
		httpClient: httpClient,
	}
}

func (c *OxfordClient) Word(appID, appKey string, lang string, word string) (*EntriesResponse, error) {
	if lang == "" {
		lang = defaultLang
	}
	urlPath := fmt.Sprintf("/entries/%s/%s", lang, word)
	const allFields = "definitions,domains,etymologies,examples,pronunciations"

	u := mustGetBaseURL()
	u.Path = path.Join(u.Path, urlPath)

	q := u.Query()
	q.Set("fields", allFields)
	q.Set("strictMatch", "false")
	u.RawQuery = q.Encode()

	resp, err := c.doGet(u, appID, appKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &EntriesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *OxfordClient) doGet(url *url.URL, appID, appKey string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("app_id", appID)
	req.Header.Set("app_key", appKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type EntriesResponse struct {
	ID       string `json:"id"`
	Metadata struct {
		Operation string `json:"operation"`
		Provider  string `json:"provider"`
		Schema    string `json:"schema"`
	} `json:"metadata"`
	Results []struct {
		ID             string `json:"id"`
		Language       string `json:"language"`
		Lexicalentries []struct {
			Entries []struct {
				Etymologies     []string `json:"etymologies"`
				Homographnumber string   `json:"homographNumber"`
				Pronunciations  []struct {
					Dialects         []string `json:"dialects"`
					Phoneticnotation string   `json:"phoneticNotation"`
					Phoneticspelling string   `json:"phoneticSpelling"`
					Audiofile        string   `json:"audioFile,omitempty"`
				} `json:"pronunciations"`
				Senses []struct {
					Definitions []string `json:"definitions"`
					Examples    []struct {
						Text string `json:"text"`
					} `json:"examples"`
					ID        string `json:"id"`
					Subsenses []struct {
						Definitions []string `json:"definitions"`
						Examples    []struct {
							Text string `json:"text"`
						} `json:"examples"`
						ID        string `json:"id"`
						Registers []struct {
							ID   string `json:"id"`
							Text string `json:"text"`
						} `json:"registers,omitempty"`
					} `json:"subsenses,omitempty"`
					Registers []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"registers,omitempty"`
				} `json:"senses"`
			} `json:"entries"`
			Language        string `json:"language"`
			Lexicalcategory struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"lexicalCategory"`
			Text string `json:"text"`
		} `json:"lexicalEntries"`
		Type string `json:"type"`
		Word string `json:"word"`
	} `json:"results"`
	Word string `json:"word"`
}
