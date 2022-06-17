// Package fastmath can be used to create simple math tasks
package fastmath

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	Addition = iota
	Subtraction
	Multiplication
	Division
	OperationsCount
	Difficulty = 10
)

// MathTask is a main struct for fastmath package
type MathTask struct {
	question string
	answer   string
}

// Generate is a method that generates single math task
func (m *MathTask) Generate() {
	rand.Seed(time.Now().UnixNano())
	a := rand.Intn(Difficulty)       //nolint:gosec // this result is not used in a secure application
	b := rand.Intn(Difficulty)       //nolint:gosec // this result is not used in a secure application
	op := rand.Intn(OperationsCount) //nolint:gosec // this result is not used in a secure application
	switch op {
	case Addition:
		m.question = fmt.Sprintf("%d + %d = ", a, b)
		m.answer = strconv.Itoa(a + b)
	case Subtraction:
		m.question = fmt.Sprintf("%d - %d = ", a, b)
		m.answer = strconv.Itoa(a - b)
	case Multiplication:
		m.question = fmt.Sprintf("%d * %d = ", a, b)
		m.answer = strconv.Itoa(a * b)
	case Division:
		m.question = fmt.Sprintf("%d / %d = ", a, b)
		m.answer = strconv.FormatFloat(float64(a)/float64(b), 'f', 1, 64) //nolint:gomnd // Standard parameters
	}
}

// GetQuestion is a getter for a question field
func (m *MathTask) GetQuestion() string {
	return m.question
}

// GetAnswer is a getter for an answer field
func (m *MathTask) GetAnswer() string {
	return m.answer
}

// SetQuestion is a setter for a question field
func (m *MathTask) SetQuestion(question string) {
	m.question = question
}

// SetAnswer is a setter for an answer field
func (m *MathTask) SetAnswer(answer string) {
	m.answer = answer
}

// SetAll is a setter for all fields
func (m *MathTask) SetAll(question, answer string) {
	m.question = question
	m.answer = answer
}
