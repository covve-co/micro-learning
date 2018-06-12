package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func SetFlash(w http.ResponseWriter, m string) {
	c := &http.Cookie{
		Name:  "flash",
		Value: encode([]byte(m)),
	}
	http.SetCookie(w, c)
}

func GetFlash(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie("flash")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return "", fmt.Errorf("no flash")
		default:
			return "", err
		}
	}

	v, err := decode(c.Value)
	if err != nil {
		return "", err
	}

	dc := &http.Cookie{Name: "flash", MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return v, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) (string, error) {
	v, err := base64.URLEncoding.DecodeString(src)
	return string(v), err
}
