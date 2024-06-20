package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

const VAULT_KUBERNETES_AUTH_METHOD_TYPE = "kubernetes"

func getAllK8sAuthMethodMountPaths(client *api.Client) ([]string, error) {
	authMethods, err := client.Sys().ListAuthWithContext(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error listing vault auth methods: %v", err)
	}

	allK8sAuthMethodMountPaths := []string{}

	for authMethodMountPath, authMethodInfo := range authMethods {
		if authMethodInfo.Type == VAULT_KUBERNETES_AUTH_METHOD_TYPE {
			allK8sAuthMethodMountPaths = append(allK8sAuthMethodMountPaths, authMethodMountPath)
		}
	}

	return allK8sAuthMethodMountPaths, nil
}
