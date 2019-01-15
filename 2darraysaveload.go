package main

import (
	"bufio"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// 10000,15000 くらいにすると差ははっきりするが、
// 数百メガバイト程度のファイルサイズになるので注意
// （現時点でも 1,000*1,000*4 = 4Mbyte以上にはなる）
const rmax = 1000
const cmax = 1000

// A は元データとなる配列。これをファイルに保存する。
var A [rmax][cmax]int32

// Bb はファイルから binary.Read で取ったデータを保存する配列
var Bb [rmax][cmax]int32

// Bg はファイルから gob.Decode で取ったデータを保存する配列
var Bg [rmax][cmax]int32

func main() {

	start := time.Now()
	fmt.Println("2D array data read and write test")
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)

	// 0. 二次元配列を乱数で作る

	rand.Seed(11111)

	start = time.Now()

	fmt.Println("set random data to 2D array A")

	for r := 0; r < rmax; r++ {
		for c := 0; c < cmax; c++ {
			A[r][c] = rand.Int31()
		}
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列サイズを測っておく
	fmt.Println("raw of A:", len(A))
	fmt.Println("column of A:", len(A[0]))
	fmt.Println("======")

	// 1. 二次元配列をバイナリファイルとして保存 binary.Write and bufio

	filebW, err := os.Create("./binaryfile.bw")
	if err != nil {
		fmt.Println("file couldn't open")
		panic(err)
	}
	defer filebW.Close()

	fmt.Println("write A to binaryfile with binary.Write and bufio")

	start = time.Now()

	// 本当はbytes で書きたかったんだけどやり方がよくわからなかったので
	// bufio で書いてしまう作戦。
	bufbW := bufio.NewWriter(filebW)

	// 書き込みバッファに書く
	err = binary.Write(bufbW, binary.LittleEndian, A)
	if err != nil {
		fmt.Println("bufio.Write failed:", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 2. 二次元配列をバイナリファイルとして保存 gob

	filegW, err := os.Create("./binaryfile.gob")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer filegW.Close()

	fmt.Println("write A to binaryfile with gob.Encode")

	start = time.Now()

	// io.Writer に向けたエンコーダーを作る。io.Writerは書き込み可として開いたファイル
	encgob := gob.NewEncoder(filegW)

	// Encodeメソッドに送り付けたい変数 A を与える。
	// この場合はファイルに対して送り付けられる。
	if err = encgob.Encode(A); err != nil {
		fmt.Println("data couldn't Encode", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	//
	fmt.Println("======")

	// 3. バイナリファイルから二次元配列 Bb に読み込み binary.Read and bufio

	filebR, err := os.Open("./binaryfile.bw")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer filebR.Close()

	fmt.Println("read from binaryfile to Bb with binary.Read and bufio")

	start = time.Now()

	// 本当はbytes で読みたかったんだけどやり方がよくわからなかったので
	// bufio で読んでしまう作戦。
	// bufR := new(bytes.Buffer)

	bufbR := bufio.NewReader(filebR)

	// binary.Read で　ファイルに紐づけたバッファから 変数に読み込む。
	// 変数はポインタで与える
	err = binary.Read(bufbR, binary.LittleEndian, &Bb)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列比較
	if A == Bb {
		fmt.Println("A is same as Bb")
	} else {
		fmt.Println("A is not same as Bb")
	}

	fmt.Printf(" A:%d %d %d %d\n", A[0][0], A[0][cmax-1], A[rmax-1][0], A[rmax-1][cmax-1])
	fmt.Printf("Bb:%d %d %d %d\n", Bb[0][0], Bb[0][cmax-1], Bb[rmax-1][0], Bb[rmax-1][cmax-1])

	// 4. バイナリファイルから二次元配列に読み込み gob

	filegR, err := os.Open("./binaryfile.gob")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer filegR.Close()

	fmt.Println("read from binaryfile to Bg with gob.Decode")

	start = time.Now()

	// io.Writer から受けるデコーダーを作る。io.Writerは開いたファイルのことになる
	decgob := gob.NewDecoder(filegR)

	// Decodeメソッドに受け取りたい変数のポインタを与える。
	// 送られてくる（ファイルの中に書き込まれている）データと同じ形、
	// サイズの変数じゃないとエラーになる。賢い。
	if err = decgob.Decode(&Bg); err != nil {
		fmt.Println("decgob.Decode failed:", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// Bg[rmax-1][cmax-1] = 111111

	// 配列比較
	if A == Bg {
		fmt.Println("A is same as Bg")
	} else {
		fmt.Println("A is not same as Bg")
	}

	fmt.Printf(" A:%d %d %d %d\n", A[0][0], A[0][cmax-1], A[rmax-1][0], A[rmax-1][cmax-1])
	fmt.Printf("Bg:%d %d %d %d\n", Bg[0][0], Bg[0][cmax-1], Bg[rmax-1][0], Bg[rmax-1][cmax-1])

}
