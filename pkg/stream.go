package pkg

import (
	"fmt"
	"strings"
)

type Callback[T any] func(T)

type Stream[T any] func(Callback[T])

func Of[T any](t T) Stream[T] {
	return func(c Callback[T]) {
		c(t)
	}
}

func From[T any](ts ...T) Stream[T] {
	return func(c Callback[T]) {
		for _, t := range ts {
			c(t)
		}
	}
}

type Mapper[T, R any] func(T) R

func Map[T, R any](s Stream[T], mapper Mapper[T, R]) Stream[R] {
	return func(c Callback[R]) {
		s(func(t T) {
			r := mapper(t)
			c(r)
		})
	}
}

type FlatMapper[T, E any] func(T) Stream[E]

func FlatMap[T, E any](s Stream[T], mapper FlatMapper[T, E]) Stream[E] {
	return func(c Callback[E]) {
		s(func(t T) {
			mapper(t)(c)
		})
	}
}

func (s Stream[T]) Filter(f func(T) bool) Stream[T] {
	return func(c Callback[T]) {
		s(func(t T) {
			if f(t) {
				c(t)
			}
		})
	}
}

func (s Stream[T]) Take(n int) Stream[T] {
	return func(callback Callback[T]) {
		s.ConsumeTillStop(func(t T) {
			if n <= 0 {
				panic("stop")
			}
			callback(t)
			n--
		})
	}
}

func (s Stream[T]) Drop(n int) Stream[T] {
	return func(callback Callback[T]) {
		s(func(t T) {
			if n > 0 {
				n--
			} else {
				callback(t)
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

func (s Stream[T]) ConsumeTillStop(callback Callback[T]) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	s(func(t T) {
		callback(t)
	})
}

func Zip[T, E, R any](s Stream[T], els []E, f func(T, E) R) Stream[R] {
	var i int
	return func(callback Callback[R]) {
		s.ConsumeTillStop(func(t T) {
			if i >= len(els) {
				panic("stop")
			}
			callback(f(t, els[i]))
			i++
		})
	}

}

func underscoreToCamelCase(str string) string {
	capitalize := func(s string) string {
		return strings.ToUpper(s[0:1]) + s[1:]
	}

	stream := func(c Callback[func(string) string]) {
		c(strings.ToLower)
		for {
			c(capitalize)
		}
	}

	return Zip(stream, strings.Split(str, "_"), func(f func(string) string, str string) string {
		return f(str)
	}).Join("")
}
