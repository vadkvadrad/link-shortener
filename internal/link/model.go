package link

import (
	"adv-go/api/internal/stat"
	"math/rand"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
	Stats []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete: SET NULL;"`
}

func NewLink(url string) *Link {
	link :=  &Link{
		Url: url,
	}
	link.GenerateHash()
	return link
}

func (l *Link) GenerateHash() {
	l.Hash = RandStringRunes(6)
}

var lettersRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range n {
		b[i] = lettersRunes[rand.Intn(len(lettersRunes))]
	}
	return string(b)
}