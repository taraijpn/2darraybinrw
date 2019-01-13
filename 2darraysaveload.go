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

// 10000,10000 くらいにすると数百メガバイト程度のファイルサイズになる
const rmax = 1000
const cmax = 1000

var a [rmax][cmax]int32
var bbw [rmax][cmax]int32
var bgob [rmax][cmax]int32

func main() {

	start := time.Now()
	fmt.Println("2D array data read and write test")
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)

	// 0. 二次元配列を乱数で作る

	rand.Seed(11111)

	start = time.Now()

	fmt.Println("set data to 2D array")

	for r := 0; r < rmax; r++ {
		for c := 0; c < cmax; c++ {
			a[r][c] = rand.Int31()
		}
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列サイズを測っておく
	fmt.Println("raw:", len(a))
	fmt.Println("column:", len(a[0]))

	// 1. 二次元配列をバイナリファイルとして保存 binary.Write and bufio

	file, err := os.Create("./binaryfile.bw")
	if err != nil {
		fmt.Println("file couldn't open")
		panic(err)
	}
	defer file.Close()

	fmt.Println("save to binaryfile with binary.Write and bufio")

	start = time.Now()

	// 本当はbytes で書きたかったんだけどやり方がよくわからなかったので
	// bufio で書いてしまう作戦。
	bufW := bufio.NewWriter(file)

	// 書き込みバッファに書く
	err = binary.Write(bufW, binary.LittleEndian, a)
	if err != nil {
		fmt.Println("bufio.Write failed:", err)
		panic(err)
	}

	// file close
	file.Close()

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 2. 二次元配列をバイナリファイルとして保存 gob

	filegob, err := os.Create("./binaryfile.gob")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer file.Close()

	fmt.Println("save to binaryfile with gob.Encode and bufio")

	start = time.Now()

	// io.Writer に向けたエンコーダーを作る。io.Writerは書き込み可として開いたファイル
	encgob := gob.NewEncoder(filegob)

	// Encodeメソッドに送り付けたい変数 a を与える。
	// この場合はファイルに対して送り付けられる。
	if err = encgob.Encode(a); err != nil {
		fmt.Println("data couldn't Encode", err)
		panic(err)
	}

	// file close
	filegob.Close()

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 3. バイナリファイルから二次元配列に読み込み binary.Read and bufio

	fileR, err := os.Open("./binaryfile.bw")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer file.Close()

	fmt.Println("load from binaryfile with binary.Read and bufio")

	start = time.Now()

	// 本当はbytes で読みたかったんだけどやり方がよくわからなかったので
	// bufio で読んでしまう作戦。
	// bufR := new(bytes.Buffer)

	bufR := bufio.NewReader(fileR)

	// binary.Read で　ファイルに紐づけたバッファから 変数に読み込む。
	// 変数はポインタで与える
	err = binary.Read(bufR, binary.LittleEndian, &bbw)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列比較
	fmt.Printf("a:%d %d %d %d\n", a[0][0], a[0][cmax-1], a[rmax-1][0], a[rmax-1][cmax-1])
	fmt.Printf("b:%d %d %d %d\n", bbw[0][0], bbw[0][cmax-1], bbw[rmax-1][0], bbw[rmax-1][cmax-1])

	if a == bbw {
		fmt.Printf("a is same as bbw\n")
	} else {
		fmt.Printf("a is not same as bbw\n")
	}

	// 4. バイナリファイルから二次元配列に読み込み gob

	filegobR, err := os.Open("./binaryfile.gob")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer file.Close()

	fmt.Println("load from binaryfile with gob")

	start = time.Now()

	// io.Writer から受けるデコーダーを作る。io.Writerは開いたファイルのことになる
	decgob := gob.NewDecoder(filegobR)

	// Decodeメソッドに受け取りたい変数のポインタを与える。
	// 送られてくるデータ（ファイル）と同じ形、サイズの変数じゃないとエラーになる。賢い。
	if err = decgob.Decode(&bgob); err != nil {
		fmt.Println("decgob.Decode failed:", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列比較
	fmt.Printf("a:%d %d %d %d\n", a[0][0], a[0][cmax-1], a[rmax-1][0], a[rmax-1][cmax-1])
	fmt.Printf("b:%d %d %d %d\n", bgob[0][0], bgob[0][cmax-1], bgob[rmax-1][0], bgob[rmax-1][cmax-1])

	if a == bgob {
		fmt.Printf("a is same as bgob\n")

	} else {
		fmt.Printf("a is not same as bgob\n")
	}

}
