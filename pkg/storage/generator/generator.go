package generator

import (
	"math/rand"
	"time"
)

type Generator struct {
	alphabet string
	length   int
}

//New - конструктор объекта генератора
func New(alphabet string, shortLength int) *Generator {
	rand.Seed(time.Now().UnixNano())
	g := &Generator{}
	g.alphabet = alphabet
	g.length = shortLength
	return g
}

// Do - генерирует строку нужной длины из заданного алфавита
func (g *Generator) Do() string {
	abc := []byte(g.alphabet)
	short := []byte{}
	for i := 1; i <= g.length; i++ {
		short = append(short, abc[rand.Intn(len(abc)-1)])
	}
	return string(short)
}
