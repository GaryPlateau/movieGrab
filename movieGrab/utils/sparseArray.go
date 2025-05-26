// 稀疏数组
package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type SparseArr struct {
	Row int `json:"row"`
	Col int `json:"col"`
	Val int `json:"val"`
}

func displayMap(maps []byte) {
	var tmps [][]int
	err := json.Unmarshal(maps, &tmps)
	if err != nil {
		fmt.Println("反序列化maps失败", err)
		return
	}
	for _, v1 := range tmps {
		for _, v2 := range v1 {
			fmt.Printf("%v\t", v2)
		}
		fmt.Println()
	}
}

func CreateMap(row int, col int) (maps []byte) {
	if row == 1 || col == 1 {
		maps = append(maps, byte(rand.Intn(100)))
		return
	}

	var tmp [][]int
	//var arr []int

	for i := 0; i < row; i++ {
		arr := make([]int, col)
		tmp = append(tmp, arr)
	}
	// 取当前时间戳
	var timeStamp = time.Now().UnixNano()
	// 构造一个rand，并使用时间戳作为他的随机种子
	r := rand.New(rand.NewSource(timeStamp))
	var tmpNum int
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			tmpNum = r.Intn(100)
			if tmpNum >= 60 {
				tmp[i][j] = tmpNum
			}
		}
	}

	maps, err := json.Marshal(tmp)
	if err != nil {
		fmt.Println("序列化map失败", err)
		return nil
	}
	return
}

func ChessMap(maps []byte) {
	displayMap(maps)
	var chessMap [][]int
	err := json.Unmarshal(maps, &chessMap)
	if err != nil {
		fmt.Println("反序列话map错误", err)
		return
	}

	var chesses []SparseArr
	for i := 0; i < len(chessMap); i++ {
		for j := 0; j < len(chessMap[i]); j++ {
			if chessMap[i][j] != 0 {
				chesses = append(chesses, SparseArr{
					Row: i,
					Col: j,
					Val: chessMap[i][j],
				})
			}
		}
	}

	//fmt.Println(chesses)
}
