package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
	"github.com/aws/aws-sdk-go-v2/service/codedeploy/types"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
)

var (
	version = "unknown"
)

func main() {
	fmt.Printf("Drone AWS CodeDeploy Plugin built from %s\n", version)

	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}

	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if vargs.Application == "" {
		vargs.Application = repo.Name
	}

	if vargs.RevisionType == "" {
		vargs.RevisionType = string(types.RevisionLocationTypeGitHub)
	}

	if vargs.AccessKey == "" {
		fmt.Println("Please provide an access key id")
		os.Exit(1)
	}

	if vargs.SecretKey == "" {
		fmt.Println("Please provide a secret access key")
		os.Exit(1)
	}

	if vargs.Region == "" {
		fmt.Println("Please provide a region")
		os.Exit(1)
	}

	if vargs.DeploymentGroup == "" {
		fmt.Println("Please provide a deployment group")
		os.Exit(1)
	}

	var location *types.RevisionLocation

	switch types.RevisionLocationType(vargs.RevisionType) {
	case types.RevisionLocationTypeGitHub:
		location = &types.RevisionLocation{
			RevisionType: types.RevisionLocationTypeGitHub,
			GitHubLocation: &types.GitHubLocation{
				CommitId:   aws.String(build.Commit),
				Repository: aws.String(repo.FullName),
			},
		}
	case types.RevisionLocationTypeS3:
		if vargs.BundleType == "" {
			fmt.Println("Please provide a bundle type")
			os.Exit(1)
		}

		if vargs.BucketName == "" {
			fmt.Println("Please provide a bucket name")
			os.Exit(1)
		}

		if vargs.BucketKey == "" {
			fmt.Println("Please provide a bucket key")
			os.Exit(1)
		}

		switch types.BundleType(vargs.BundleType) {
		case types.BundleTypeTar:
		case types.BundleTypeTarGZip: // tgz
		case types.BundleTypeZip:
		default:
			fmt.Println("Invalid bundle type")
			os.Exit(1)
		}

		s3location := &types.S3Location{
			BundleType: types.BundleType(vargs.BundleType),
			Bucket:     aws.String(vargs.BucketName),
			Key:        aws.String(vargs.BucketKey),
		}

		if vargs.BucketEtag != "" {
			s3location.ETag = aws.String(vargs.BucketEtag)
		}

		if vargs.BucketVersion != "" {
			s3location.Version = aws.String(vargs.BucketVersion)
		}

		location = &types.RevisionLocation{
			RevisionType: types.RevisionLocationTypeS3,
			S3Location:   s3location,
		}
	default:
		fmt.Println("Invalid revision type")
		os.Exit(1)
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(vargs.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			vargs.AccessKey,
			vargs.SecretKey,
			"",
		)),
	)
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		os.Exit(1)
	}

	svc := codedeploy.NewFromConfig(cfg)

	_, err = svc.CreateDeployment(ctx,
		&codedeploy.CreateDeploymentInput{
			ApplicationName:               aws.String(vargs.Application),
			DeploymentConfigName:          aws.String(vargs.DeploymentConfig),
			DeploymentGroupName:           aws.String(vargs.DeploymentGroup),
			Description:                   aws.String(vargs.Description),
			IgnoreApplicationStopFailures: vargs.IgnoreStopFailures,
			Revision:                      location,
		},
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Successfully deployed")
}
