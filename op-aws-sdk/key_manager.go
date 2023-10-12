package op_aws_sdk

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"os"
)

const (
	OP_NODE_P2P_SEQUENCER_KEY = "OP_NODE_P2P_SEQUENCER_KEY"
	OP_BATCHER_SIGN_KEY       = "OP_BATCHER_SIGN_KEY"
	OP_PROPOSER_SIGN_KEY      = "OP_PROPOSER_SIGN_KEY"
)

func Key_manager(context context.Context, ctx *cli.Context, keyName string) error {

	aws_key_id := ""
	aws_key_region := ""
	key_flag_name := ""
	log.Info("Key manager ", "keyName", keyName)
	switch keyName {
	case OP_NODE_P2P_SEQUENCER_KEY:
		aws_key_id = "AWS_P2P_SEQUENCER_KEY_ID"
		aws_key_region = "AWS_P2P_SEQUENCER_KEY_REGION"
		key_flag_name = "p2p.sequencer.key"
	case OP_BATCHER_SIGN_KEY:
		aws_key_id = "AWS_OP_BATCHER_SIGN_KEY_ID"
		aws_key_region = "AWS_OP_BATCHER_SIGN_KEY_REGION"
		key_flag_name = "private-key"
	case OP_PROPOSER_SIGN_KEY:
		aws_key_id = "AWS_OP_PROPOSER_SIGN_KEY_ID"
		aws_key_region = "AWS_OP_PROPOSER_SIGN_KEY_REGION"
		key_flag_name = "private-key"
	default:
		log.Error("Key manager ", "error keyName", keyName)
		return nil
	}

	return load(context, ctx, aws_key_id, aws_key_region, key_flag_name)
}
func load(context context.Context, ctx *cli.Context, aws_key_region string, aws_key_id string, key_flag_name string) error {
	aws_p2p_sequencer_key_id := os.Getenv(aws_key_id)
	aws_p2p_sequencer_key_region := os.Getenv(aws_key_region)
	if aws_p2p_sequencer_key_id != "" || aws_p2p_sequencer_key_region != "" {
		log.Info("Key manager ", "aws_key_region", aws_p2p_sequencer_key_region, "aws_key_id", aws_p2p_sequencer_key_id)
		config, err := config.LoadDefaultConfig(context, config.WithRegion(aws_p2p_sequencer_key_region))
		if err != nil {
			log.Error("Key manager load key config from aws", "error", err)
			return err
		}
		secretManager := secretsmanager.NewFromConfig(config)
		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(aws_p2p_sequencer_key_id),
			VersionStage: aws.String("AWSCURRENT"),
		}
		result, err := secretManager.GetSecretValue(context, input)
		if err != nil {
			log.Error("Key manager key value from aws", "error", err)
			return err
		}
		log.Info("Key manager load key ", "key", *result.SecretString)
		ctx.Set(key_flag_name, *result.SecretString)
	}
	return nil
}
