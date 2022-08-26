package wxcommon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const wxBaseURL = "https://qyapi.weixin.qq.com/cgi-bin"

type TokenHandler func() (string, error)

type WorkClient struct {
	GetAccessToken TokenHandler
}

func (wx *WorkClient) GetJSON(api string, resp WorkWxResp) error {
	r, err := http.Get(wxBaseURL + api)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("wxbizhttp:%d(%s)", r.StatusCode, r.Status)
	}

	if err = json.NewDecoder(r.Body).Decode(resp); err != nil {
		return err
	}

	return resp.Err()
}

func (wx *WorkClient) GetRespWithToken(api string, resp WorkWxResp, args ...interface{}) error {
	token, err := wx.GetAccessToken()
	if err != nil {
		return err
	}

	var argsWithToken = []interface{}{token}
	argsWithToken = append(argsWithToken, args...)

	return wx.GetJSON(fmt.Sprintf(api, argsWithToken...), resp)
}

func (wx *WorkClient) PostJSON(api string, req interface{}, resp WorkWxResp) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := http.Post(wxBaseURL+api, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("wxbizhttp:%d(%s)", r.StatusCode, r.Status)
	}

	if err = json.NewDecoder(r.Body).Decode(resp); err != nil {
		return err
	}

	return resp.Err()
}

func (wx *WorkClient) PostMedia(api string, media *MediaToUpload, resp WorkWxResp) error {
	buf := bytes.NewBuffer(nil)
	mv := multipart.NewWriter(buf)
	wr, err := mv.CreateFormFile("media", media.filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, media.r)
	if err != nil {
		return err
	}

	r, err := http.Post(wxBaseURL+api, mv.FormDataContentType(), buf)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("wxbizhttp:%d(%s)", r.StatusCode, r.Status)
	}

	if err = json.NewDecoder(r.Body).Decode(resp); err != nil {
		return err
	}

	return resp.Err()
}

func (wx *WorkClient) PostRespWithToken(api string, req interface{}, resp WorkWxResp, args ...interface{}) error {
	token, err := wx.GetAccessToken()
	if err != nil {
		return err
	}

	var argsWithToken = []interface{}{token}
	argsWithToken = append(argsWithToken, args...)

	return wx.PostJSON(fmt.Sprintf(api, argsWithToken...), req, resp)
}
