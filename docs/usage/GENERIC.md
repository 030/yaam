# generic

## Configuration

~/.yaam/conf/repositories/generic.yaml

```bash
---
allowedRepos:
  - something
```

## Upload

```bash
curl -X POST -u hello:world http://yaam.some-domain/generic/something/world4.iso --data-binary @/home/${USER}/Downloads/ubuntu-22.04.1-desktop-amd64.iso
```

### Troubleshooting

```bash
413 Request Entity Too Large
```

add:

```bash
data:
  proxy-body-size: 5G
```

and restart the controller pod.

Verify in the `/etc/nginx/nginx.conf` file that the `client_max_body_size` has
been increased to 5G.

## Download

```bash
curl -u hello:world http://yaam.some-domain/generic/something/world6.iso -o /tmp/world6.iso
```
