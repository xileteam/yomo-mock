# ys5

#### Build
```sh
sh build.sh
```

#### Run
```sh
./build/ys5_proxy

YOMO_SFN_PORT=12001 ./build/ys5_crawler

curl -L -x socks5://localhost:8888 http://yomo.run
```
