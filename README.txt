ip.chuhlomin.com is a service for finding information about IP addresses.

It uses:
  * GeoLite2 databases for ASN and GeoIP lookups,
  * whois.iana.org for Whois queries.

Available endpoints:
  /help - this page
  / - index page, redirects to /{ip}, where {ip} is your IP address
  /{ip} - returns information about the IP address: ASN and GeoIP
  /{ip}/whois - returns the Whois information for the IP address
  /{ip}/{mask} - displays the IP in binary format, visualizing the mask

Example usage:
  curl -L https://ip.chuhlomin.com/
  curl https://ip.chuhlomin.com/1.1.1.1
  curl https://ip.chuhlomin.com/1.1.1.1/whois
  curl https://ip.chuhlomin.com/192.168.0.0/24
  curl -s https://ip.chuhlomin.com/1.1.1.1.json | jq -r '.asn.number'

Version: 1.0.0
Source code: https://github.com/chuhlomin/ip
Author: Konstantin Chukhlomin
License: MIT

---

Known alternatives:
  https://ifconfig.co
  https://ipinfo.io
  https://whatismyipaddress.com
