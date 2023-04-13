package pkg

import (
	"fmt"
	"strings"
)

type Consumer[T any] func(t T)

type Stream[T any] func(consumer Consumer[T])

func From[T any](ts ...T) Stream[T] {
	return func(consumer Consumer[T]) {
		for _, t := range ts {
			consumer(t)
		}
	}
}

type Mapper[T, R any] func(T) R

func Map[T, R any](s Stream[T], mapper Mapper[T, R]) Stream[R] {
	return func(consumer Consumer[R]) {
		s(func(t T) {
			r := mapper(t)
			consumer(r)
		})
	}
}

type FlatMapper[T, E any] func(T) Stream[E]

func FlatMap[T, E any](s Stream[T], mapper FlatMapper[T, E]) Stream[E] {
	return func(consumer Consumer[E]) {
		s(func(t T) {
			mapper(t)(consumer)
		})
	}
}

func (s Stream[T]) Filter(f func(T) bool) Stream[T] {
	return func(consumer Consumer[T]) {
		s(func(t T) {
			if f(t) {
				consumer(t)
			}
		})
	}
}

func (s Stream[T]) Take(n int) Stream[T] {
	return func(consumer Consumer[T]) {
		s.ConsumeTillStop(func(t T) {
			if n <= 0 {
				panic("stop")
			}
			consumer(t)
			n--
		})
	}
}

func (s Stream[T]) Drop(n int) Stream[T] {
	return func(consumer Consumer[T]) {
		s(func(t T) {
			if n > 0 {
				n--
			} else {
				consumer(t)
			}
		})
	}
}

func (s Stream[T]) Join(sep string) string {
	var ss []string
	s(func(t T) {
		ss = append(ss, String(t))
	})
	return strings.Join(ss, sep)
}

func (s Stream[T]) Array() []T {
	var ss []T
	s(func(t T) {
		ss = append(ss, t)
	})
	return ss
}

func String[T any](t T) string {
	return fmt.Sprintf("%v", t)
}

func (s Stream[T]) ConsumeTillStop(consumer Consumer[T]) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	s(func(t T) {
		consumer(t)
	})
}

func Zip[T, E, R any](s Stream[T], els []E, f func(T, E) R) Stream[R] {
	var i int
	return func(consumer Consumer[R]) {
		s.ConsumeTillStop(func(t T) {
			if i >= len(els) {
				panic("stop")
			}
			consumer(f(t, els[i]))
			i++
		})
	}

}

func underscoreToCamelCase(str string) string {
	capitalize := func(s string) string {
		return strings.ToUpper(s[0:1]) + s[1:]
	}

	stream := func(c Consumer[func(string) string]) {
		c(strings.ToLower)
		for {
			c(capitalize)
		}
	}

	zipped := Zip(stream, strings.Split(str, "_"), func(f func(string) string, str string) string {
		return f(str)
	})
	return zipped.Join("")
}
