package streaming

type Config struct {
	URL       string `koanf:"url"`
	ClusterID string `koanf:"clusterid"`
}
