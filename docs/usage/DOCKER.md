# docker

[![dockeri.co](https://dockeri.co/image/utrecht/yaam)](https://hub.docker.com/r/utrecht/yaam)

```bash
docker run \
  -v /home/${USER}/.yaam/conf:/opt/yaam/.yaam/conf \
  -v /home/${USER}/.yaam/repositories:/opt/yaam/.yaam/repositories \
  -e YAAM_USER=hello \
  -e YAAM_PASS=world \
  -p 25213:25213 \
  -it utrecht/yaam:0.2.1
```
