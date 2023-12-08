package op_aws_sdk

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

const (
	OP_NODE_P2P_SEQUENCER_KEY = "OP_NODE_P2P_SEQUENCER_KEY"
	OP_BATCHER_SIGN_KEY       = "OP_BATCHER_SIGN_KEY"
	OP_PROPOSER_SIGN_KEY      = "OP_PROPOSER_SIGN_KEY"
	AWS_KEY_JSON_NAME         = "pk"
)

func KeyManager(context context.Context, ctx *cli.Context, keyType string) error {
	secretName := ""
	awsRegion := ""
	flagName := ""
	log.Info("Key manager ", "keyType", keyType)
	switch keyType {
	case OP_NODE_P2P_SEQUENCER_KEY:
		secretName = "OP_NODE_AWS_P2P_SECRET_NAME"
		awsRegion = "OP_NODE_AWS_P2P_SECRET_REGION"
		flagName = "p2p.sequencer.key"
	case OP_BATCHER_SIGN_KEY:
		secretName = "OP_BATCHER_AWS_SECRET_NAME"
		awsRegion = "OP_BATCHER_AWS_SECRET_REGION"
		flagName = "private-key"
	case OP_PROPOSER_SIGN_KEY:
		secretName = "OP_PROPOSER_AWS_SECRET_NAME"
		awsRegion = "OP_PROPOSER_AWS_SECRET_REGION"
		flagName = "private-key"
	default:
		log.Error("Key manager ", "unknown keyType", keyType)
		panic("Key manager unknown key type")
	}
	return load(context, ctx, awsRegion, secretName, flagName)
}
func load(context context.Context, ctx *cli.Context, awsRegion string, secretName string, flagName string) error {
	name := os.Getenv(secretName)
	region := os.Getenv(awsRegion)
	if name != "" && region != "" {
		loadKeyConfig, err := config.LoadDefaultConfig(context, config.WithRegion(region))
		if err != nil {
			log.Error("Key manager load key config from aws", "error", err)
			return err
		}
		secretManager := secretsmanager.NewFromConfig(loadKeyConfig)
		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(name),
			VersionStage: aws.String("AWSCURRENT"),
		}
		result, err := secretManager.GetSecretValue(context, input)
		if err != nil {
			log.Error("Key manager key value from aws", "error", err)
			return err
		}
		resultMap := make(map[string]string)
		secretBytes := []byte(*result.SecretString)
		err = json.Unmarshal(secretBytes, &resultMap)
		if err != nil {
			return err
		}
		key, ok := resultMap[AWS_KEY_JSON_NAME]
		if !ok {
			log.Error("Key manager load key does not exist")
			return errors.New("Key manager load key does not exist")
		}
		log.Info("Key manager load key is success")
		ctx.Set(flagName, key)
	} else {
		log.Info("Key manager is skipped")
	}
	return nil
}
