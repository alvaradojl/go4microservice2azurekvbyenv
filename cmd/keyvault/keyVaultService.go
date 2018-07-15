package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/opentracing-contrib/go-stdLib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
)

var errSecretNotFound = errors.New("secret not found")
var errNoAccessToken = errors.New("no access token")
var errNoKeyPath = errors.New("no key path found")

// AccessToken represents the response from Azure AD auth endpoint
type AccessToken struct {
	Token_type     string
	Expires_in     string
	Ext_expires_in string
	Expires_on     string
	Not_before     string
	Resource       string
	Access_token   string
}

var accessToken string

type requestBody struct {
	KeyPath string
}

type keyVaultResponse struct {
	Value       string
	ContentType string
	Id          string
	Attributes  keyvaultResponseAttributes
}

type keyvaultResponseAttributes struct {
	Enabled       bool
	Created       int64
	Updated       int64
	RecoveryLevel string
}

type keyVaultConfig struct {
	KeyVaultAuthEndpoint        string
	KeyVaultAuthClientID        string
	KeyVaultAuthClientSecret    string
	KeyVaultSvcBaseEndpoint     string
	KeyVaultAuthGrantType       string
	KeyVaultAuthResource        string
	KeyVaultAPIVersionParameter string
}

func getKeyVaultConfig() keyVaultConfig {
	var result = keyVaultConfig{}

	viper.SetEnvPrefix("relayservices")
	viper.BindEnv("keyvaultauthendpoint")
	viper.BindEnv("keyvaultauthclientid")
	viper.BindEnv("keyvaultauthclientsecret")
	viper.BindEnv("keyvaultsvcbaseendpoint")

	//get config from environment
	result.KeyVaultAuthEndpoint = viper.Get("keyvaultauthendpoint").(string)
	result.KeyVaultAuthClientID = viper.Get("keyvaultauthclientid").(string)
	result.KeyVaultAuthClientSecret = viper.Get("keyvaultauthclientsecret").(string)
	result.KeyVaultSvcBaseEndpoint = viper.Get("keyvaultsvcbaseendpoint").(string)

	result.KeyVaultAuthGrantType = "client_credentials"
	result.KeyVaultAuthResource = "https://vault.azure.net"
	result.KeyVaultAPIVersionParameter = "?api-version=2016-10-01"

	return result
}

func getSecret(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	span, _ := opentracing.StartSpanFromContext(ctx, "getSecret")
	defer span.Finish()

	var rb requestBody

	//only react to POST method, everything else should be 404
	if r.Method != http.MethodPost {
		http.Error(w, "method not found", http.StatusNotFound)
	}

	//keypath is provided in the payload
	//keypath should look like secrets/<name>/<version>
	err1 := json.NewDecoder(r.Body).Decode(&rb)

	if err1 != nil {
		http.Error(w, "unable to read payload", http.StatusBadRequest)
	}

	if len(rb.KeyPath) <= 0 {
		http.Error(w, "no keypath provided", http.StatusBadRequest)
	}

	infoCtx(ctx, "retrieving access token...")

	kvconf := getKeyVaultConfig()

	url := kvconf.KeyVaultSvcBaseEndpoint + rb.KeyPath + kvconf.KeyVaultAPIVersionParameter

	var err2 error
	if len(accessToken) == 0 {
		accessToken, err2 = getAccessToken(ctx)

		if err2 != nil {
			errorCtx(ctx, err2)
			http.Error(w, "unable to retrieve secret (err2)", http.StatusInternalServerError)
		}
	}

	bearer := "Bearer " + accessToken

	req, err3 := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	if err3 != nil {
		errorCtx(ctx, err3)
		http.Error(w, "unable to retrieve secret (err3)", http.StatusInternalServerError)
	}

	req = req.WithContext(ctx)

	req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req)

	defer ht.Finish()

	resp, err4 := httpClient.Do(req)

	if err4 != nil {
		errorCtx(ctx, err4)
		http.Error(w, "unable to retrieve secret (err4)", http.StatusInternalServerError)
	}

	if resp.StatusCode == http.StatusOK {

		var kvResponse keyVaultResponse

		err5 := json.NewDecoder(resp.Body).Decode(&kvResponse)

		if err5 != nil {
			errorCtx(ctx, err5)
			http.Error(w, "unable to retrieve secret (err5)", http.StatusInternalServerError)
		}

		_, _ = w.Write([]byte(kvResponse.Value))
	} else if resp.StatusCode == http.StatusUnauthorized {
		//TODO: refactor with circuit braker to get access token and try again multiple times
		http.Error(w, "unable to retrieve secret (unathorized)", http.StatusInternalServerError)

	}

}

func getAccessToken(ctx context.Context) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "getAccessToken")
	defer span.Finish()

	kvconf := getKeyVaultConfig()

	//	infoCtx(ctx, "login to keyvault access token auth : "+kvconf.KeyVaultAuthEndpoint)

	form := url.Values{}
	form.Add("grant_type", kvconf.KeyVaultAuthGrantType)
	form.Add("client_id", kvconf.KeyVaultAuthClientID)
	form.Add("client_secret", kvconf.KeyVaultAuthClientSecret)
	form.Add("resource", kvconf.KeyVaultAuthResource)

	req, err1 := http.NewRequest(http.MethodPost, kvconf.KeyVaultAuthEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err1 != nil {
		errorCtx(ctx, err1)
		return "", err1
	}

	req = req.WithContext(ctx)

	req, ht := nethttp.TraceRequest(opentracing.GlobalTracer(), req)

	defer ht.Finish()

	token := new(AccessToken)

	resp, err2 := httpClient.Do(req)
	if err2 != nil {
		errorCtx(ctx, err2)
		return "", err2
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//	bodyString := string(bodyBytes)

	//	infoCtx(ctx, "response body: "+bodyString)

	err3 := json.Unmarshal(bodyBytes, token)

	if err3 != nil {
		errorCtx(ctx, err3)
		return "", err3
	}

	//	infoCtx(ctx, "access token found: "+token.Access_token)

	return token.Access_token, nil

}
