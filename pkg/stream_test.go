package pkg

import (
	"reflect"
	"strconv"
	"testing"
)

func TestFrom(t *testing.T) {
	s := From(0, 1, 2)
	joined := s.Join("_")
	expected := "0_1_2"
	if joined != expected {
		t.Errorf("expected: %s, actual: %s", expected, joined)
	}

	s2 := From(0)
	joined = s2.Join("_")
	expected = "0"
	if joined != expected {
		t.Errorf("expected: %s, actual: %s", expected, joined)
	}
}

func TestMapFunction(t *testing.T) {
	s := func(c Consumer[int]) {
		for i := 0; i < 10; i++ {
			c(i)
		}
	}
	i := Map(s, func(t int) string {
		return strconv.Itoa(t)
	})
	joined := i.Join("_")
	expected := "0_1_2_3_4_5_6_7_8_9"
	if joined != expected {
		t.Errorf("Map test failed: expected '%s', got '%s'", expected, joined)
	}
}

func TestFlatMap(t *testing.T) {
	s := From(1, 2, 3)
	f := func(t int) Stream[int] {
		return From(t, t*10, t*20)
	}
	expected := "1 10 20 2 20 40 3 30 60"
	result := FlatMap(s, f).Join(" ")
	if result != expected {
		t.Errorf("FlatMap test failed: expected '%s', got '%s'", expected, result)
	}
}

func TestFilter(t *testing.T) {

	stream := From(1, 2, 3, 4, 5)

	filter := func(t int) bool {
		return t%2 == 0
	}

	filteredStream := stream.Filter(filter)
	expected := []int{2, 4}

	var actual []int
	consumer := func(t int) {
		actual = append(actual, t)
	}

	filteredStream(consumer)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Filter test failed: expected %v, but got %v", expected, actual)
	}
}

func TestTake(t *testing.T) {
	input := From(1, 2, 3, 4, 5)
	output := input.Take(3)

	var result []int
	output(func(i int) {
		result = append(result, i)
	})

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Take test failed: expected %v, but got %v", expected, result)
	}
}

func TestDrop(t *testing.T) {

	stream := From(1, 2, 3, 4, 5)
	dropped := stream.Drop(2)
	expected := []int{3, 4, 5}

	var actual []int
	dropped(func(i int) {
		actual = append(actual, i)
	})

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Drop test failed: expected %v, but got %v", expected, actual)
	}
}

func TestArray(t *testing.T) {
	s := From(1, 2, 3)
	arr := s.Array()
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("Array test failed: expected %v, but got %v", expected, arr)
	}
}

func TestZipFunction(t *testing.T) {
	s := func(c Consumer[int]) {
		for i := 0; i < 10; i++ {
			c(i)
		}
	}
	i := []string{"Zero", "One", "Two"}

	zipped := Zip(s, i, func(t int, e string) string {
		return strconv.Itoa(t) + "_" + e
	})

	join := zipped.Join(",")
	expected := "0_Zero,1_One,2_Two"
	if join != expected {
		t.Errorf("expected: %s, actual: %s", expected, join)
	}
}

func TestUnderscoreToCamelCase(t *testing.T) {

	actual := underscoreToCamelCase("hello_world")
	expected := "helloWorld"

	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	actual = underscoreToCamelCase("hello_world_test")
	expected = "helloWorldTest"

	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}

	actual = underscoreToCamelCase("hello")
	expected = "hello"

	if expected != actual {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}
