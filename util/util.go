package util

import (
	"bytes"
	"github.com/google/go-querystring/query"
	"strings"
	"github.com/andream16/go-amazon-itemlookup/configuration"
	"github.com/andream16/go-amazon-itemlookup/model"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func BuildQuery(q model.LightQueryModel, c configuration.Configuration, secr string) string {
	o := orderQueryModelAttributesByBytes(q)
	p := prepareForSignature(o, c)
	s := signRequest(p, secr)
	e := escapeSpecialCharacters(s)
	return getFinalRequest(e, o, c)
}

func orderQueryModelAttributesByBytes(q model.LightQueryModel) string {
	qs := QueryModelToQueryString(q)
	qa := strings.Split(qs, "&")
	b := stringSliceToBytesSlice(qa)
	s := mergeSort(b)
	var lastIndex = len(s)-1
	var f string
	for i := range s {
		f += string(s[i][:])
		if i != lastIndex {
			f += "&"
		}
	}
	return f
}

func prepareForSignature(o string, c configuration.Configuration) string {
	return c.Remote.Verb + "\n" +
		c.Remote.Endpoint + "\n" +
		"/" + c.Remote.Service + "/" + c.Remote.Format + "\n" +
		o
}

func signRequest(r string, secr string) string {
	key := []byte(secr)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(r))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func escapeSpecialCharacters(s string) string {
	m := strings.Replace(s, "+", "%2B", -1)
	return strings.Replace(m, "=", "%3D", -1)
}

func getFinalRequest(signature , object string, c configuration.Configuration) string {
	return "http://" + c.Remote.Endpoint + "/" + c.Remote.Service + "/" + c.Remote.Format + "?" + object + "&Signature=" + signature
}

func QueryModelToQueryString (q interface{}) string {
	v, _ := query.Values(q)
	return v.Encode()
}

func stringSliceToBytesSlice(s []string) [][]byte {
	b := [][]byte{}
	for _, k := range s {
		b = append(b, []byte(k))
	}
	return b
}

func mergeSort(a [][]byte) [][]byte {

	if len(a) <= 1 {
		return a
	}

	left := make([][]byte, 0)
	right := make([][]byte, 0)
	m := len(a) / 2

	for i, x := range a {
		switch {
		case i < m:
			left = append(left, x)
		case i >= m:
			right = append(right, x)
		}
	}

	left = mergeSort(left)
	right = mergeSort(right)

	return merge(left, right)
}

func merge(left, right [][]byte) [][]byte {

	results := make([][]byte, 0)

	for len(left) > 0 || len(right) > 0 {
		if len(left) > 0 && len(right) > 0 {
			if bytes.Compare(left[0], right[0]) == -1 {
				results = append(results, left[0])
				left = left[1:]
			} else {
				results = append(results, right[0])
				right = right[1:]
			}
		} else if len(left) > 0 {
			results = append(results, left[0])
			left = left[1:]
		} else if len(right) > 0 {
			results = append(results, right[0])
			right = right[1:]
		}
	}

	return results
}