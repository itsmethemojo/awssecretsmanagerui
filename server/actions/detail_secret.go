package actions

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func CheckNameInList(listName []string, name string) bool {
	for i := range listName {
		if listName[i] == name {
			return true
		}
	}

	return false
}

func GetSecretByARN(region, arn string) (secretsmanager.DescribeSecretOutput, error) {
  
	config, configErr := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if configErr != nil {
		return secretsmanager.DescribeSecretOutput{}, configErr
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(arn),
	}
	result, err := svc.DescribeSecret(context.TODO(), input)

	filterNames := GetFilterNames()
	if len(filterNames) > 0 && !CheckNameInList(filterNames, *result.Name) {
		return secretsmanager.DescribeSecretOutput{}, errors.New("Can't get secret")
	}

	if err != nil {
		return secretsmanager.DescribeSecretOutput{}, err
	}
	if result == nil {
		return secretsmanager.DescribeSecretOutput{}, errors.New("Can't get secret")
	}
	return *result, nil
}
