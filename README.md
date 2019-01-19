# 2次元配列をバイナリデータとして読み書きする

2次元の整数配列をバイナリデータとしてファイルに読み書きする方法をgoの練習がてら試してみたもの。
いろいろすっきりしない書き方をしている気がする。

Package binary のドキュメントには
> Clients that require high-performance serialization, especially for large data structures, should look at more advanced solutions such as the encoding/gob package or protocol buffers.
https://golang.org/pkg/encoding/binary/

とか書いてあるし、gobのほうがいいのかな？
実際のところ、gob.Encodeはbinary.Writeより遅いけど、gob.Decodeはbinary.Readより早い。

binary.Read が遅いのはデータの一つ一つについてエンディアンなどを確認してるから。
だとすれば、読み書きする環境が同じ、と限定して、バイナリ列として扱えればそれでよいのでは？
というアイディアを採用すると、読み書きともに桁違いに早くなる。詳細は以下URL参照。
https://bleu48.hatenablog.com/entry/2019/01/19/185959

[10000][15000]int32 で試してみた結果はこちら。
動作OSは Windows 10 Home (64bit)で、ドライブはSSDになっている。
この結果を信じるなら file.Read() が早すぎると思う。
600MByte のバイナリファイルを 0.119秒で読み込んでるんだけど、
SATA/600(600Mbyte/sec) の規格を超えているような…？？

```
$ ./2darraysaveload.exe > result.txt
2D array data read and write test
0s
set random data to 2D array A
3.5344796s
raw of A: 10000
column of A: 15000
======
write A to binaryfile with binary.Write and bufio
5.12895s
set random data to 2D array B
A is not same as B
A:1179337617 1132468563 2094486022 368894420
B:611703184 1569009249 1390812853 720412150
read from binaryfile to B with binary.Read and bufio
3.0586963s
A is same as B
A:1179337617 1132468563 2094486022 368894420
B:1179337617 1132468563 2094486022 368894420
======
write A to binary gob file with gob.Encode()
6.0728727s
set random data to 2D array B
A is not same as B
A:1179337617 1132468563 2094486022 368894420
B:1760913819 1904829126 1299570517 1622806500
read from binary gob file to B with gob.Decode
1.9258039s
A is same as B
A:1179337617 1132468563 2094486022 368894420
B:1179337617 1132468563 2094486022 368894420
======
write A to binaryfile with file.Write()
1.7527747s
set random data to 2D array B
A is not same as B
A:1179337617 1132468563 2094486022 368894420
B:824297579 1169316836 212125022 1884237302
binaryfilesize: 600000000
arraysize:      600000000
read from binaryfile to B with file.Read and unsafe pointer casting
118.9047ms
A is same as B
A:1179337617 1132468563 2094486022 368894420
B:1179337617 1132468563 2094486022 368894420

$ md5sum.exe binaryfile.* >> result.txt
1e05fdaee78bf097b457ccb0846a3e98 *binaryfile.bw
1e05fdaee78bf097b457ccb0846a3e98 *binaryfile.fw
20da06f7ce6ead41541e128ec74c5aff *binaryfile.gob
```

