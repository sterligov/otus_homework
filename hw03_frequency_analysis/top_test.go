package hw03_frequency_analysis //nolint:golint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
			require.Subset(t, expected, Top10(text))
		} else {
			expected := []string{"он", "и", "а", "что", "ты", "не", "если", "-", "то", "Кристофер"}
			require.ElementsMatch(t, expected, Top10(text))
		}
	})

	t.Run("text with runes from many languages", func(t *testing.T) {
		text := `
		Go — язык общего назначения с широкими широкими широкими возможностями и понятным синтаксисом
		The Go programming language is an open source project to make programmers more productive.

		123 123 123 123 123 test9 ,* . @#$

		Go - 編譯型 种静态强类型 編譯型 并发型 并具有垃圾回收功能的编程语言 編譯型
		Go ist eine kompilierbare Programmiersprache, die Nebenläufigkeit unterstützt und über eine automatische Speicherbereinigung verfügt.
		Entwickelt wurde Go von Mitarbeitern des Unternehmens Google Inc
		Go es un lenguaje de programación de fuente abierta diseñado para construir software simple, rápido y confiable.

		Go is expressive, concise, clean, and efficient. Its concurrency mechanisms make it easy to write programs that
		get the most out of multicore and networked machines, while its novel type system enables flexible and modular
		program construction. Go compiles quickly to machine code yet has the convenience of garbage collection and
		the power of run-time reflection. It's a fast, statically typed, compiled language that feels like a dynamically typed,
		interpreted language.
		`
		actual := Top10(text)
		expected := []string{"the", "and", "широкими", "go", "eine", "language", "編譯型", "de", "der", "to", "of", "is", "make", "its", "that", "a", "typed"}

		require.Equal(t, 10, len(actual))
		require.Subset(t, expected, actual)
	})

	t.Run("less then 10 different words", func(t *testing.T) {
		text := `The 'Go-Go' programming programming programming programming programming programming language is programming language.`
		expected := []string{"the", "go-go", "programming", "language", "is"}

		actual := Top10(text)

		require.Equal(t, len(expected), len(actual))
		require.Subset(t, expected, actual)
	})
}

func TestMapWordFrequency(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		actual := MapWordFrequency("")
		expected := map[string]int{}
		require.Equal(t, expected, actual)
	})

	t.Run("not empty string", func(t *testing.T) {
		text := `
		Словом efficient зовут efficient    Винни-Пух
		Go Go expressive    минутку efficient`

		expected := map[string]int{
			"словом":     1,
			"efficient":  3,
			"зовут":      1,
			"go":         2,
			"винни-пух":  1,
			"expressive": 1,
			"минутку":    1,
		}

		actual := MapWordFrequency(text)

		require.Equal(t, expected, actual)
	})
}

func TestMapFrequencyWords(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		wordFrequency, maxFrequency := MapFrequencyWords(nil)
		expected := map[int][]string{}

		require.Equal(t, 0, maxFrequency)
		require.Equal(t, len(expected), len(wordFrequency))
	})

	t.Run("not empty map", func(t *testing.T) {
		words := map[string]int{
			"словом":     1,
			"efficient":  3,
			"зовут":      1,
			"go":         2,
			"the":        2,
			"винни-пух":  1,
			"expressive": 1,
			"минутку":    1,
		}

		expected := map[int][]string{
			1: {"словом", "зовут", "винни-пух", "expressive", "минутку"},
			2: {"go", "the"},
			3: {"efficient"},
		}

		wordFrequency, maxFrequency := MapFrequencyWords(words)

		require.Equal(t, 3, maxFrequency)
		require.Equal(t, len(expected), len(wordFrequency))

		for k, v := range expected {
			require.Contains(t, wordFrequency, k)
			require.Equal(t, len(v), len(wordFrequency[k]))
			require.Subset(t, v, wordFrequency[k])
		}
	})
}
