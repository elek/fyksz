package k8s

type K8s struct {
	Env         Env         `cmd:"" help:"transformer to add/replace environment variables in deployments"`
	AsConfigMap AsConfigMap `cmd:"" help:"standalone helper to convert raw file to K8s ConfigMap"`
}
