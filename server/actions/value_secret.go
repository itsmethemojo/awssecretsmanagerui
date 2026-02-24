package actions

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecretValueByARN(region, arn string) (secretsmanager.GetSecretValueOutput, error) {
	config, configErr := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if configErr != nil {
		return secretsmanager.GetSecretValueOutput{}, configErr
	}
	
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(arn),
	}
	result, err := svc.GetSecretValue(context.TODO(), input)

	filterNames := GetFilterNames()
	if len(filterNames) > 0 && !CheckNameInList(filterNames, *result.Name) {
		return secretsmanager.GetSecretValueOutput{}, errors.New("Can't get secret")
	}

	if err != nil {
		return secretsmanager.GetSecretValueOutput{}, err
	}
	if result == nil {
		return secretsmanager.GetSecretValueOutput{}, errors.New("Can't get secret")
	}
	return *result, nil
}

func UpdateSecretValue(region string, request secretsmanager.PutSecretValueInput) (secretsmanager.GetSecretValueOutput, error) {
	config, configErr := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if configErr != nil {
		return secretsmanager.GetSecretValueOutput{}, configErr
	}
	
	svc := secretsmanager.NewFromConfig(config)

	if _, err := GetSecretByARN(region, *request.SecretId); err != nil {
		return secretsmanager.GetSecretValueOutput{}, errors.New("Can't update secret")
	}

	_, err := svc.PutSecretValue(context.TODO(), &request)
	if err != nil {
		return secretsmanager.GetSecretValueOutput{}, err
	}
	return GetSecretValueByARN(region, *request.SecretId)
}

func UpdateSecretValueBinary(region string, arn string, binaryVaue []byte) (secretsmanager.GetSecretValueOutput, error) {
	config, configErr := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if configErr != nil {
		return secretsmanager.GetSecretValueOutput{}, configErr
	}
	
	svc := secretsmanager.NewFromConfig(config)

	if _, err := GetSecretByARN(region, arn); err != nil {
		return secretsmanager.GetSecretValueOutput{}, errors.New("Can't upload secret")
	}

	request := secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(arn),
		SecretBinary: binaryVaue,
	}
	_, err := svc.PutSecretValue(context.TODO(),&request)
	if err != nil {
		return secretsmanager.GetSecretValueOutput{}, err
	}
	return GetSecretValueByARN(region, arn)
}
