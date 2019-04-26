package main

import (
	"encoding/base64"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/pcantea/portainer-ecr/portainer"
)

type Portainer struct {
	RegistryIds []string
	client      *portainer.Client
}

func main() {

	ticker := time.NewTicker(6 * time.Hour)

	go func() {
		for ; true; <-ticker.C {
			client, err := portainer.NewClient(
				os.Getenv("PORTAINER_URL"),
				os.Getenv("PORTAINER_USER"),
				os.Getenv("PORTAINER_PASS"),
				nil,
			)
			if err != nil {
				log.Fatal("Could not initialize Portainer client: ", err)
			}

			p := Portainer{
				client: client,
			}

			if ids, ok := os.LookupEnv("AWS_ECR_REGISTRY_IDS"); ok && ids != "" {
				p.RegistryIds = strings.Split(ids, ",")
			}
			p.getECRToken(newEcrClient())
		}
	}()

	select {}
}

func (p *Portainer) getECRToken(svc ecriface.ECRAPI) {
	request := &ecr.GetAuthorizationTokenInput{}
	if len(p.RegistryIds) > 0 {
		request = &ecr.GetAuthorizationTokenInput{RegistryIds: aws.StringSlice(p.RegistryIds)}
	}

	resp, err := svc.GetAuthorizationToken(request)
	if err != nil {
		log.Fatal("Could not retrieve ECR Authorization token", err)
	}
	if len(resp.AuthorizationData) < 1 {
		log.Fatal("Request did not return properly formated authorization data")
	}

	for _, data := range resp.AuthorizationData {
		p.updateToken(data)
	}
}

func (p *Portainer) updateToken(data *ecr.AuthorizationData) {
	bytes, err := base64.StdEncoding.DecodeString(*data.AuthorizationToken)
	if err != nil {
		log.Fatal(err)
	}
	token := string(bytes[:len(bytes)])

	authTokens := strings.Split(token, ":")
	if len(authTokens) != 2 {
		log.Fatal("Authorization token does not contain correct data")
		return
	}

	registryURL, err := url.Parse(*data.ProxyEndpoint)
	ecrUser := authTokens[0]
	ecrPass := authTokens[1]
	ecrURL := registryURL.Host

	registries := p.client.GetRegistries()
	for _, registry := range registries {
		if registry.URL == ecrURL {
			registry.Username = ecrUser
			registry.Password = ecrPass
		}
		p.client.UpdateRegistry(&registry)

	}
}

func newEcrClient() *ecr.ECR {
	roleArn, ok := os.LookupEnv("AWS_ROLE_ARN")
	if ok {
		return ecr.New(
			session.New(
				aws.NewConfig().WithCredentials(
					stscreds.NewCredentials(session.New(), roleArn),
				),
			),
		)
	}

	return ecr.New(session.New())
}
