package service

func MakeDiscoveryAddress(serviceName string) string {
	return "discovery:///" + serviceName
}

func MakeEtcdAddress(serviceName string) string {
	return "etcd://localhost:2379/" + serviceName + "/"
}
