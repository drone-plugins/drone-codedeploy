package main

type Params struct {
	AccessKey          string `json:"access_key"`
	SecretKey          string `json:"secret_key"`
	Region             string `json:"region"`
	Application        string `json:"application"`
	DeploymentGroup    string `json:"deployment_group"`
	DeploymentConfig   string `json:"deployment_config"`
	Description        string `json:"description"`
	RevisionType       string `json:"revision_type"`
	BundleType         string `json:"bundle_type"`
	BucketName         string `json:"bucket_name"`
	BucketKey          string `json:"bucket_key"`
	BucketEtag         string `json:"bucket_etag"`
	BucketVersion      string `json:"bucket_version"`
	IgnoreStopFailures bool   `json:"ignore_stop_failures"`
}
