package hw02unpackstring

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	inputIndex := 0
	resultIndex := 0
	backslash := '\\'
	input := []rune(str)
	result := make([]rune, 0, 200)
	inputLen := len(input) - 1

	// validate length
	if utf8.RuneCountInString(string(input)) == 0 {
		return "", nil
	}

	// starts with digit
	if unicode.IsDigit(input[0]) {
		return "", ErrInvalidString
	}

	for inputIndex <= inputLen {
		r := input[inputIndex]

		inputNextIndex := nextIndex(&input, inputIndex)
		nextR := input[inputNextIndex]
		// finds number
		if inputLen != inputNextIndex && unicode.IsDigit(r) && unicode.IsDigit(nextR) {
			fmt.Println("finds number")
			return "", ErrInvalidString
		}

		// escaped symbol is valid
		isValidForEscape := unicode.IsDigit(nextR) || nextR == backslash
		if r == backslash && isValidForEscape {
			result = append(result, nextR)
			resultIndex++
			inputIndex += 2
			continue
		}

		// escaped symbol is not valid
		if r == backslash && !isValidForEscape {
			fmt.Println("escaped symbol is not valid")
			return "", ErrInvalidString
		}

		// repeatedly add rune
		isDigit := unicode.IsDigit(r)
		if isDigit {
			shift := appendRune(&result, result[previousIndex(&result, resultIndex)], r)
			resultIndex += shift
			inputIndex++
			continue
		}

		// remove current rune
		vl, err := strconv.Atoi(string(nextR))
		if err == nil && vl == 0 {
			inputIndex++
			continue
		}

		// simple append
		result = append(result, r)
		resultIndex++
		inputIndex++
	}

	return string(result), nil
}

func appendRune(res *[]rune, rn rune, count rune) int {
	value, err := strconv.Atoi(string(count))
	if err != nil {
		return 0
	}

	for range value - 1 {
		*res = append(*res, rn)
	}
	return value
}

func nextIndex(sl *[]rune, index int) int {
	slLen := math.Max(0, float64(len(*sl)-1))
	nextIndex := index + 1
	return int(math.Min(float64(slLen), float64(nextIndex)))
}

func previousIndex(sl *[]rune, index int) int {
	slLen := math.Max(0, float64(len(*sl)-1))
	previousIndex := index - 1
	return int(math.Min(float64(slLen), math.Max(0, float64(previousIndex))))
}
