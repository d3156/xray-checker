docker run --rm \
  -p 2112:2112 \
  -e SUBSCRIPTION_URL='https://raw.githubusercontent.com/igareck/vpn-configs-for-russia/refs/heads/main/WHITE-CIDR-RU-checked.txt' \
  -e EXPORT_TOKEN='test-token' \
  -e EXPORT_BASE64='false' \
  -e EXPORT_MAX_LATENCY_MS='5000' \
  -e EXPORT_MAX_NODES='5' \
  xray-checker:local