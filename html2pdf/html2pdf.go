package html2pdf

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deltrinos/tpl21/log"
	"io/ioutil"
	"net/http"
)

type WkPdfRequest struct {
	Contents string `json:"contents"`
}

type PdfServer struct {
	Url    string
	Req    WkPdfRequest
	Result []byte
}

func NewPdfServer(url string) *PdfServer {
	return &PdfServer{Url: url}
}

func (p *PdfServer) ConvertString(html string) error {
	p.Req.Contents = base64.StdEncoding.EncodeToString([]byte(html))
	return p.makeRequest()
}

func (p *PdfServer) makeRequest() error {
	p.Result = []byte{}

	payloadBytes, err := json.Marshal(p.Req)
	if err != nil {
		log.Error().Msgf("failed to json.Marshal html %v", err)
		return err
	}

	payload := bytes.NewReader(payloadBytes)
	request, err := http.NewRequest("POST", p.Url, payload)
	if err != nil {
		log.Error().Msgf("failed to create NewRequest %v", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error().Msgf("failed to make client request %v", err)
		return err
	}
	defer resp.Body.Close()

	p.Result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("failed to read response body %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status is not success %v", resp.StatusCode))
	}
	return nil
}

func (p *PdfServer) ConvertFile(path string) error {
	res, err := readFileInBase64(path)
	if err != nil {
		log.Error().Msgf("failed to readFileInBase64 %v", err)
		return err
	}
	return p.ConvertString(res)
}

func readFileInBase64(path string) (string, error) {
	var res string
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return res, err
	}
	res = base64.StdEncoding.EncodeToString(f)
	return res, nil
}

func (p *PdfServer) WriteToFile(path string) error {
	if len(p.Result) <= 0 {
		return errors.New("response result is empty")
	}
	return ioutil.WriteFile(path, p.Result, 0644)
}
