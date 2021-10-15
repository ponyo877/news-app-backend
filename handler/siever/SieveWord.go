package siever

import (
	_ "fmt"
	"regexp"
	"strings"

	"github.com/tomoemon/text_normalizer"
)

func normalizeTitle(title string) (normalizedTitle string) {
	normalizer := text_normalizer.NewTextNormalizer(
		text_normalizer.AlphabetToHankaku,
		text_normalizer.KanaToHiragana)
	normalizedTitle = normalizer.Replace(title)
	rep := regexp.MustCompile("[^0-9a-zA-Zぁ-んァ-ヶ一-龠ー]+")
	normalizedTitle = rep.ReplaceAllString(normalizedTitle, "@")
	// fmt.Println(normalizedTitle)
	return
}

func ContainsNGWord(title string) bool {
	normalizedTitle := normalizeTitle(title)
	for _, NGword := range NGwords {
		// fmt.Println(NGword)
		if strings.Contains(normalizedTitle, NGword) {
			return true
		}
	}
	return false
}
