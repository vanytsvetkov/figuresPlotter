package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

type Screen struct {
	rows uint16
	cols uint16
	Xpx  uint16
	Ypx  uint16
}

func initScreen() *Screen {
	// clear terminal window
	fmt.Printf("\033[2J")

	// define terminal's size
	//cols, _ := strconv.Atoi(os.Getenv("COLUMNS"))
	//rows, _ := strconv.Atoi(os.Getenv("LINES"))
	//return &Screen{rows, cols}

	screen := new(Screen)

	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(1074295912),
		uintptr(unsafe.Pointer(screen)),
	)
	return screen
}

func (screen Screen) clear() {
	// return cursor to top left
	fmt.Printf("\033[H")
}

func (screen Screen) printT(matrix [][]int) {
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("%2d", matrix[i][j])
		}
		fmt.Println()
	}
}

func (screen Screen) printM(matrix [][]string) {
	for i := range matrix {
		for j := range matrix[i] {
			fmt.Printf("\033[32;1m")
			//fmt.Printf("%s", color.GreenString(fmt.Sprintf("%2d", matrix[i][j])))
			fmt.Printf("%2s", matrix[i][j])
			fmt.Printf("\033[0m")
		}
		fmt.Println()
	}
}

func genTriangle(rows int, cols int, line int) [][]int {
	matrix := make([][]int, rows)

	for i := range matrix {
		matrix[i] = make([]int, cols)
		start := len(matrix[i])/2 - i
		end := len(matrix[i])/2 + i
		for j := start; j <= end; j++ {
			if i <= line {
				matrix[i][j] = 1
			}
		}
	}

	return matrix
}

func genMatrix(limit, rows, cols int, matrix [][]string) [][]string {
	for k := 0; k <= limit; k++ {
		n := rand.Intn(5)
		j := rand.Intn(cols)
		if n == 0 {
			for i := range matrix {
				matrix[i][j] = ""
			}
			continue
		}
		for i := range matrix {
			// Generate a random integer between 32 and 126 (inclusive)
			// ASCII codes for printable characters range from 32 to 126
			randInt := rand.Intn(95) + 32
			matrix[i][j] = string(rune(randInt))
		}
	}
	return matrix
}

func interruptHandler(screen Screen) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interruptChannel

		screen.clear()
		os.Exit(0)
	}()
}

func figures(screen Screen) {

	fmt.Println(screen)
	time.Sleep(2 * time.Second)

	cols := int(screen.cols) / 2
	rows := int(screen.rows) - 1

	/*	Triangle
		for i := 0; i <= cols/2; i++ {
			screen.clear()
			matrix := genTriangle(rows, cols, i)
			screen.printT(matrix)
			time.Sleep(1 * time.Second)
		}
	*/

	/*	Matrix	*/
	matrix := make([][]string, rows)
	for i := range matrix {
		matrix[i] = make([]string, cols)
	}
	for {
		screen.clear()
		matrix = genMatrix(15, rows, cols, matrix)
		screen.printM(matrix)
		time.Sleep(350 * time.Millisecond)
	}

}

func main() {
	screen := initScreen()

	interruptHandler(*screen)
	figures(*screen)
	select {}
}
