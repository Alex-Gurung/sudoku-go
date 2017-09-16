package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//"strings"
//"errors"
//"math/rand"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func import_file(s string) string {
	dat, err := ioutil.ReadFile(s)
	check(err)
	return string(dat)
}
func make_matrix(s string) [][]string {
	/*
		...11.4..3..2...
		should become
		[. . . 1]
		[1 . 4 .]
		[. 3 . .]
		[2 . . .]

	*/
	var side_length = int(math.Sqrt(float64(len(s)))) //Find the side length of the board (NOT THE SQUARE LENGTH)
	matrix := make([][]string, side_length)           //Create an empty 2 dimensional slice to store the strings
	for i := 0; i < side_length; i++ {                //For loops to create the matrix
		for x := 0; x < side_length; x++ {
			matrix[i] = append(matrix[i], string(s[i*side_length+x])) //Append current 'character' (it's a string) to the current row being iterated over
		}
	}
	return matrix //Return the matrix
}
func print_matrix(m [][]string) { //Simple method to print out board, to be nicely formatted later
	for _, value := range m {
		fmt.Println(value)
	}
	fmt.Println()
}
func print_array(a []string) {
	fmt.Println(strings.Join(a, ","))
}

var standard_global = []string{}

var global_options = []string{"1", "2", "3", "4"}

func set_global_options(m []string) { //setter for global variable
	global_options = m
}

/*
* POTENTIAL ALTERNATIVE TO CONSIDER:
* maybe instead of checking a move should check board? that seems backwards/more
* inefficient though since it's essentially doing this for every space
 */
func move_legal(matrix [][]string, row_pos int, col_pos int, choice string) bool { //evaluates a move based on the board
	/** CHECK EMPTY **/
	if matrix[row_pos][col_pos] != "." {
		return false
	}
	/** CHECK ROW **/
	for r := 0; r < len(matrix); r++ {
		if matrix[r][col_pos] == choice {
			return false
		}
	}
	/** CHECK COLUMN **/
	for c := 0; c < len(matrix); c++ {
		if matrix[row_pos][c] == choice {
			return false
		}
	}
	/** CHECK SQUARE **/
	/* FIND SQUARE BOUNDS */
	square_length := int(math.Sqrt(float64(len(matrix)))) //9x9 has a square length of 3, 4x4 - 2, etc.
	r_lower_bound := (row_pos / square_length) * square_length
	c_lower_bound := (col_pos / square_length) * square_length

	for r := r_lower_bound; r < r_lower_bound+square_length; r++ {
		for c := c_lower_bound; c < c_lower_bound+square_length; c++ {
			if matrix[r][c] == choice {
				return false
			}
		}
	}
	return true //If all conditions are met, i.e. the row, columns, and square are clear, the move is legal
}

func board_complete(matrix [][]string) (r int, c int) {
	for r := 0; r < len(matrix); r++ {
		for c := 0; c < len(matrix); c++ {
			if matrix[r][c] == "." {
				return r, c
			}
		}
	}
	return -1, -1
}

func choose_move(matrix [][]string) (r int, c int, s string) { //returns row_pos, col_pos, choice
	square_length := len(matrix)
	for r := 0; r < len(matrix); r++ {
		for c := 0; c < len(matrix); c++ {
			if matrix[r][c] == "." {
				choice := strconv.Itoa(rand.Intn(square_length) + 1)
				return r, c, choice
			}
		}
	}
	return -1, -1, "no"
}

// func smarter_choose_move(matrix [][]string) (r int, c int, s string) {
// 	for r:= 0; r < len(matrix); r++ {

// 	}
// }

func choose_move_for_space(matrix [][]string, r int, c int, start int) (s string, index int) { //returns row_pos, col_pos, choice
	for i := start; i < len(global_options); i++ {
		value := global_options[i]

		if move_legal(matrix, r, c, value) {
			return value, i
		}

	}
	return "no", -1
}

func solve_board(matrix [][]string) [][]string {
	// print_matrix(matrix)
	r, c := board_complete(matrix)
	if r == -1 { //If the board is filled, return the board (it should never be wrong as moves are checked before they are made)
		return matrix
	}
	//Else there are still places needing filled on the board
	//r, c is the space of the first empty square (not the best algorithm but for now)

	set_global_options(standard_global) //Reset the options

	choice, ind := choose_move_for_space(matrix, r, c, 0)
	for ind != -1 {
		// fmt.Println(ind, ",", choice)
		matrix[r][c] = choice
		attempt := solve_board(matrix)
		if len(attempt) == 0 {
			matrix[r][c] = "."
			choice, ind = choose_move_for_space(matrix, r, c, ind+1)

		} else {
			return attempt
		}
	}

	return make([][]string, 0)
}

func main() {
	//fmt.Printf("Hello, world.\n") //test print
	file_data := import_file("tests/another_99.txt") //Import the board from one of the files in tests/
	// fmt.Printf(file_data)                          //Print out the raw board from the file
	problems := strings.Split(file_data, "\n")
	var matrix [][]string
	// start_time := time.Now()
	var max_s float64
	max_s = 0
	var total_elapsed float64
	total_elapsed = 0
	var num_problems float64
	num_problems = 1
	for i, problem := range problems {
		num_problems = num_problems + float64(i)
		fmt.Println("Problem:", i)
		matrix = make_matrix(problem) //Construct a matrix from the raw board
		print_matrix(matrix)                         //Print out that matrix

		standard_global = []string{}
		for i := 1; i < len(matrix)+1; i++ {
			standard_global = append(standard_global, strconv.Itoa(i))
		}
		//HARDCODED FOR NOW
		// options := []string{"1", "2", "3", "4"}
		// set_global_options(options)

		// fmt.Println(strconv.FormatBool(move_legal(matrix, 1, 3, "2")))
		start := time.Now()
		solved := solve_board(matrix)
		t := time.Now()
		elapsed := t.Sub(start)
		in_s := elapsed.Seconds()
		if max_s < in_s {
			max_s = in_s
		}
		total_elapsed = total_elapsed + in_s
		fmt.Println(elapsed)

		print_matrix(solved)
	}
	// end_time := time.Now()
	// elapsed := end_time.Sub(start_time)
	fmt.Println("Total:", total_elapsed)
	fmt.Println("Average:", total_elapsed/num_problems)
	fmt.Println("Max:", max_s)


}
