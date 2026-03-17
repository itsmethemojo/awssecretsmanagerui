package actions

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

var filterNames []string

func hasNameFilter() bool {
	names := os.Getenv("FILTER_NAMES")
	if names == "" {
		return false
	}
	return true
}

func hasTagFilter() bool {
	tag := os.Getenv("FILTER_TAG")
	if tag == "" || !strings.Contains(tag, "="){
		return false
	}
	return true
}

func GetFilterNames() []string {
	if filterNames != nil {
		return filterNames
	}

	names := os.Getenv("FILTER_NAMES")
	if names == "" {
		filterNames = []string{}
	} else {
		filterNames = strings.Split(names, ",")
	}

	return filterNames
}

func getTagFilter() types.Filter {
	//TODO this tiny part can be unit tested
	//TODO we can also check again if it contains =
	tag := os.Getenv("FILTER_TAG")
	tagKeyAndValue := strings.SplitN(tag, "=", 2)
	key := tagKeyAndValue[0]
	//TODO should add a second filter that checks also the value of the tag
	//value := tagKeyAndValue[1]
	nameFilters := types.Filter{
		Key:   types.FilterNameStringTypeTagKey,
		Values: []string{key},
	}

	return nameFilters
}

func getNamesFilter(listFilterNames []string) types.Filter {
	nameFilters := types.Filter{
		Key:    types.FilterNameStringTypeName,
		Values: listFilterNames,
	}

	return nameFilters
}

func GetListSecrets(region string) ([]types.SecretListEntry, error) {
	config, configErr := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if configErr != nil {
		return []types.SecretListEntry{}, configErr
	}
	
	svc := secretsmanager.NewFromConfig(config)

	maxResult := int32(100)

	index := 0
	var token *string
	secrets := []types.SecretListEntry{}
	for ; token != nil || index == 0; index++ {
		result, err := GetAPageSecrets(svc, token, maxResult)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, result.SecretList...)
		token = result.NextToken
	}
	return secrets, nil
}

func GetAPageSecrets(svc *secretsmanager.Client, token *string, maxResult int32) (*secretsmanager.ListSecretsOutput, error) {
	input := &secretsmanager.ListSecretsInput{
		MaxResults: &maxResult,
		NextToken:  token,
	}

	if hasNameFilter() {
		input.Filters = []types.Filter{
			getNamesFilter(GetFilterNames()),
		}
	} else if hasTagFilter() {
		input.Filters = []types.Filter{
			getTagFilter(),
		}
	}

	result, err := svc.ListSecrets(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result, nil
}
