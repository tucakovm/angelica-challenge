package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	totalScore := 0

	lips := findLipsFromFile("angelica.txt")

	pupils := findPupilsFromFile("angelica.txt")
	totalScore = computeScoreFromFile("angelica.txt", lips, pupils)

	fmt.Println("Total Polkadot Score:", totalScore)
}

func computePolkadotScore(runes []rune, lips [][2]int, pupils [][]rune) int {
	score := 0
	totalPupilChars := 0
	insideLips := 0
	outsideLips := 0

	for i := 0; i < len(pupils); i += 2 {
		first := pupils[i]
		second := pupils[i+1]
		totalPupilChars = totalPupilChars + len(first) + len(second)
	}

	for i, r := range runes {
		if r == 'O' {
			inside := false
			var lipRange [2]int

			for _, lip := range lips {
				if i >= lip[0] && i <= lip[1] {
					inside = true
					lipRange = lip
					break
				}
			}

			if inside {
				insideLips++
				fmt.Printf(
					"O at index %d -> INSIDE lips [%d-%d]\n",
					i, lipRange[0], lipRange[1],
				)
			} else {
				outsideLips++
				fmt.Printf(
					"O at index %d -> OUTSIDE lips\n",
					i,
				)
			}
		}
	}

	score = outsideLips + insideLips*totalPupilChars
	return score
}

func findLips(runes []rune) [][2]int {
	var segments [][2]int
	n := len(runes)
	i := 0

	for i < n {
		for i < n && unicode.IsSpace(runes[i]) {
			i++
		}
		if i >= n {
			break
		}

		start := i

		for i < n && !unicode.IsSpace(runes[i]) && !isBanned(runes[i]) {
			i++
		}

		end := i - 1
		length := end - start + 1

		if length >= 2 {
			beforeOk := start == 0 || unicode.IsSpace(runes[start-1])
			afterOk := i == n || unicode.IsSpace(runes[i])
			startOk := !isBanned(runes[start])
			endOk := !isBanned(runes[end])

			if beforeOk && afterOk && startOk && endOk {
				segments = append(segments, [2]int{start, end})
			}
		}

		if i == start {
			i++
		}
	}

	return segments
}

func findPupilCharsAll(runes []rune) [][]rune {
	var result [][]rune
	length := len(runes)

	type occ struct {
		start int
		end   int
	}

	patterns := make(map[string][]occ)

	for start := 0; start < length; start++ {
		if unicode.IsSpace(runes[start]) || isBannedExceptSpace(runes[start]) {
			continue
		}

		end := start
		for end+1 < length &&
			!unicode.IsSpace(runes[end+1]) &&
			!isBannedExceptSpace(runes[end+1]) {
			end++
		}

		segLen := end - start + 1
		if segLen < 2 {
			start = end
			continue
		}

		leftOk := start == 0 || (unicode.IsSpace(runes[start-1]) && !isBannedExceptSpace(runes[start-1]))
		rightOk := end == length-1 || (unicode.IsSpace(runes[end+1]) && !isBannedExceptSpace(runes[end+1]))

		if !leftOk || !rightOk {
			start = end
			continue
		}

		seg := string(runes[start : end+1])
		patterns[seg] = append(patterns[seg], occ{start, end})

		start = end
	}

	fmt.Println("\nFound pupil pairs: ")

	pairCounter := 0
	totalLen := 0

	for seg, occs := range patterns {
		if len(occs) == 2 {
			pairCounter++
			segLen := len(seg)

			fmt.Printf(
				"Pair %d: %q | len=%d)\n",
				pairCounter,
				seg,
				segLen,
			)

			result = append(result, []rune(seg))
			result = append(result, []rune(seg))

			totalLen += segLen * 2
		}
	}

	fmt.Printf(
		"\nTotal pupil pairs: %d\nTotal pupil characters (sum): %d\n\n",
		pairCounter,
		totalLen,
	)

	return result
}

// Helper funcs

func findLipsFromFile(filename string) [][2]int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file for lips:", err)
		return nil
	}
	defer file.Close()

	var lips [][2]int
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := []rune(scanner.Text())
		segments := findLips(line)

		for _, seg := range segments {
			start := seg[0]
			end := seg[1]
			value := string(line[start : end+1])
			fmt.Printf("Line %d, Lip: Index [%d-%d] -> %q\n",
				lineNum, start, end, value)
		}

		lips = append(lips, segments...)
	}

	return lips
}

func findPupilsFromFile(filename string) [][]rune {
	file, _ := os.Open(filename)
	defer file.Close()

	var allRunes []rune
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		allRunes = append(allRunes, []rune(scanner.Text())...)
		allRunes = append(allRunes, '\n')
	}

	return findPupilCharsAll(allRunes)
}

func computeScoreFromFile(filename string, lips [][2]int, pupils [][]rune) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file for score:", err)
		return 0
	}
	defer file.Close()

	totalScore := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []rune(scanner.Text())
		totalScore += computePolkadotScore(line, lips, pupils)
	}

	return totalScore
}

func isBanned(c rune) bool {
	return c == ',' || c == '`' || c == '\'' || c == '-' || c == 'O' || unicode.IsSpace(c)
}

func isBannedExceptSpace(c rune) bool {
	return c == ',' || c == '`' || c == '\'' || c == '-' || c == 'O'
}
