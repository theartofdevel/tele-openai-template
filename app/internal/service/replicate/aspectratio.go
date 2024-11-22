package replicate

import (
	"errors"
)

// AspectRatio — тип с перечислением возможных аспектов.
type AspectRatio int

const (
	AspectRatio21x9 AspectRatio = iota
	AspectRatio16x9
	AspectRatio3x2
	AspectRatio4x3
	AspectRatio5x4
	AspectRatio1x1
	AspectRatio4x5
	AspectRatio3x4
	AspectRatio2x3
	AspectRatio9x16
	AspectRatio9x21
)

// mapping — явное сопоставление значений AspectRatio и их строкового представления.
var mapping = map[AspectRatio]string{
	AspectRatio21x9: "21:9",
	AspectRatio16x9: "16:9",
	AspectRatio3x2:  "3:2",
	AspectRatio4x3:  "4:3",
	AspectRatio5x4:  "5:4",
	AspectRatio1x1:  "1:1",
	AspectRatio4x5:  "4:5",
	AspectRatio3x4:  "3:4",
	AspectRatio2x3:  "2:3",
	AspectRatio9x16: "9:16",
	AspectRatio9x21: "9:21",
}

// reverseMapping — обратное сопоставление для быстрого поиска по строке.
var reverseMapping = func() map[string]AspectRatio {
	rm := make(map[string]AspectRatio)
	for key, value := range mapping {
		rm[value] = key
	}
	return rm
}()

// String возвращает строковое представление AspectRatio.
func (ar AspectRatio) String() string {
	if str, ok := mapping[ar]; ok {
		return str
	}

	return "unknown"
}

// NewAspectRatio преобразует строку в значение AspectRatio с валидацией.
func NewAspectRatio(s string) (AspectRatio, error) {
	if ar, ok := reverseMapping[s]; ok {
		return ar, nil
	}

	return -1, errors.New("invalid aspect ratio string")
}
