package service

func MakeDiscoveryAddress(serviceName string) string {
	return "discovery:///" + serviceName
}
