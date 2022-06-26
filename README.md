# yomo-mock

这个repo是用来为yomo@v2.0接口进行原型设计实验的，所以不使用quic和y3，而用tcp和json实现，重点是用户接口的一致性

## Build Zipper

```sh
sh build.sh
```

## Run Zipper

```sh
./build/yomo_zipper
```

## Examples

1. [Datagram: random noise](noise/README.md)
2. [Stream: socks5 proxy](ys5/README.md)
