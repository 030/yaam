# NPM

npx create-react-app .

echo -n 'hello:world' | openssl base64

vim ~/.npmrc

registry=http://localhost:25213/npm/3rdparty-npm/
always-auth=true
_auth=aGVsbG86d29ybGQ=

mv ~/.npm/cache{,2}

cd test/npm/demo
rm -r node_modules

npm i -d

validate json using json_pp
for f in $(find /tmp/yaam/testi2/repositories/npm/3rdparty-npm/ -name *.tmp); do cat $f | json_pp; done
for f in $(find /tmp/yaam/testi2/repositories/npm/3rdparty-npm/ -name *.tmp); do du -h $f; done|grep M


ENOENT: no such file or directory
enoent x no such file or directory
try remove the package-lock.json
and run npm i -d again