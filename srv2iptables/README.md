# srv2iptables

srv2iptables builds a iptables/ip6tables chain from the content of a
DNS srv record.

## Usage

1. Set up SRV records with the hosts you want to include:

```
_kube      SRV     0 0 0 kw1
_kube      SRV     0 0 0 kw2
_kube      SRV     0 0 0 kw3
```

2. Set up a cronjob on the target host to update the chain periodically.

```crontab

7 * * * * PATH=/bin:/sbin /usr/local/bin/srv2iptables \
    --chain=FromKube--srv=_kube.example.com
```

3. Point iptables at your this chain for the port you care about:

```
iptables -I INPUT 1 -p tcp -m tcp --dport 9100 -j FromKube
```

# FAQ

## Why not use an ipset?

One of the machines I needed this is running some _ancient_ software
that doesn't support them.

Once everything is modern, we can use an ipset and point firewalld at
it: `firewall-cmd --add-rich-rule='rule source ipset=denylist drop'`

## Why not do incremental updates?

Patches welcome to remove the race condition between deleting and
re-adding.
