# NPM

Create a `.npmrc` file in the directory of a particular NPM project:

```bash
registry=http://localhost:25213/npm/3rdparty-npm/
always-auth=true
_auth=aGVsbG86d29ybGQ=
cache=/tmp/some-yaam-repo/npm/cache20220914120431999
```

Note: the `_auth` key should be populated with the output of:
`echo -n 'someuser:somepass' | openssl base64`.

```bash
npm i
```
