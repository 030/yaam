# docker

```bash
docker run -v /home/${USER}/.yaam/conf:/opt/yaam/.yaam/conf -v /home/${USER}/.yaam/repositories:/opt/yaam/.yaam/repositories -e YAAM_USER=hello -e YAAM_PASS=world -p 25213:25213 -it yaam
```
