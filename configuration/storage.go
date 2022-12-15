package configuration

type EndpointConfig struct {
	URL           string `yaml:"endpointUrl"`
	PartitionID   string `yaml:"partitionId"`
	SigningName   string `yaml:"signingName"`
	SigningMethod string `yaml:"signingMethod"`
	SigningRegion string `yaml:"signingRegion"`
}

type AccessConfig struct {
	AccessKeyID string `yaml:"accessKeyId"`
	SecretKey   string `yaml:"secretKey"`
	Token       string `yaml:"token"`
}

type BucketConfig struct {
	Name   string `yaml:"name"`
	Folder string `yaml:"folder"`
	CDNUrl string `yaml:"cdnUrl"`
}

type S3Config struct {
	Endpoint EndpointConfig `yaml:"endpoint"`
	Access   AccessConfig   `yaml:"access"`
	Bucket   BucketConfig   `yaml:"bucket"`
}

type LocalStorageConfig struct {
	StoragePath string `yaml:"storage_path"` //nolint:tagliatelle // legacy
}

type StorageConfig struct {
	Type  string             `yaml:"type"`
	Local LocalStorageConfig `yaml:"local"`
	S3    S3Config           `yaml:"s3"`
}
