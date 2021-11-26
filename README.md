# ip

```bash
docker build -t ghcr.io/chuhlomin/geolite2:latest -f Dockerfile.GeoLite2 .
docker push ghcr.io/chuhlomin/geolite2:latest
```

```bash
curl -L https://ip.chuhlomin.com
# -L flag means "follow redirects", request to root will redirect to /<your-ip>

# serves ASN & GeoIP info for given IP
curl https://ip.chuhlomin.com/1.1.1.1

curl -s https://ip.chuhlomin.com/1.1.1.1.json | jq -r '.asn.number'
# -s flag means "silent", to hide progress meter
```

```bash
https://ip.chuhlomin.com/1.1.1.1/whois
```

https://ifconfig.co
