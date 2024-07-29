package policies.helloexternal

# default to a closed system (deny by default)
default allowed = false


allowed {
	response := http.send({"method": "get", "url": "http://app:8888/external"})
	response.status_code == 200
}