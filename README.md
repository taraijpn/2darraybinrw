# 2次元配列をバイナリデータとして読み書きする
goの練習がてら試してみたけど、いろいろすっきりしない書き方をしている気がする。

gob.Encodeはbinary.Writeより遅いけど、gob.Decodeはbinary.Readより早い。
[10000][15000]Int32 で試してみた結果はこちら。

```
> go run .\2darraysaveload.go
2D array data read and write test
0s
set data to 2D array
3.5912164s
raw: 10000
column: 15000
save to binaryfile with binary.Write and bufio
5.3062375s
save to binaryfile with gob.Encode and bufio
6.4260289s
load from binaryfile with binary.Read and bufio
3.3522479s
a:1179337617 1132468563 2094486022 368894420
b:1179337617 1132468563 2094486022 368894420
a is same as bbw
load from binaryfile with gob
1.9995575s
a:1179337617 1132468563 2094486022 368894420
b:1179337617 1132468563 2094486022 368894420
a is same as bgob

> dir
Mode                LastWriteTime         Length Name
----                -------------         ------ ----
-a----       2019/01/14      1:25           4811 2darraysaveload.go
-a----       2019/01/14      1:26      600000000 binaryfile.bw
-a----       2019/01/14      1:26      749441712 binaryfile.gob
```

ファイルサイズでかいな！？

まあ Package binary のドキュメントにも 

> Clients that require high-performance serialization, especially for large data structures, should look at more advanced solutions such as the encoding/gob package or protocol buffers.

https://golang.org/pkg/encoding/binary/

って書いてあるし、gobのほうがいいのかも。
