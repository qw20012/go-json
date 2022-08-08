# go-basic

- Package basic provides some utility functions.
- Package str provides some utilities to manipulate string.

## Usage

### basic.IsNil

Identify whether the any type is nil.
```
	var emptyArray [1]int
	if basic.IsNil(emptyArray) {
		t.Fatalf("GetOrCreate with emtpy slice failed")
	}
```
### basic.NewIfEmpty

Make sure any type is created. Create by reflect if it is not there.
```
	var emptyAnyMap map[string]any
	fromEmptyAnyMap := basic.NewIfEmpty(emptyAnyMap)
	if fromEmptyAnyMap == nil {
		t.Fatalf("GetOrCreate with emtpy map failed")
	}
	fromEmptyAnyMap["key"] = 1
```
### str.Empty

Represents the emptry string.
```
	if str.Empty ！= "" {
		t.Fatal("IsEmpty failed")
	}
```
### str.IsEmpty

Identify whether the source string is empty.
```
	if IsEmpty("abc") {
		t.Fatal("IsEmpty failed " + "abc")
	}
```
### str.IsNotEmpty

Identify whether the source string is empty.
```
	if !IsNotEmpty("abc") {
		t.Fatal("IsNotEmpty failed" + "abc")
	}
```
### str.Contact

Contact the sources from any type.
```
	twoDiffType := Contact("abc", 1)
	if twoDiffType != "abc1" {
		t.Fatal("TestContact failed " + "abc, 1")
	}
```
### str.From

Convert to string from any type.
```
	f := From(1.123)
	if f != "1.123" {
		t.Fatal("From failed " + "1.123")
	}
```
### str.Format

Format source string that instead given name in curly brackets by given value.
```
	diffTypeValue := Format("abc {name}", "name", 1)
	if diffTypeValue != "abc 1" {
		t.Fatal("TestFormat failed " + "abc, 1 ")
	}
```
### str.Formats

Format source string by calling Format functon. See also Format.
```
	strFromMap := Formats("{a}{ b }c", diffTypeValue)
	if strFromMap != "Dog1c" {
		t.Fatal("TestFormats failed " + "Dog1c")
	}
```

## Contributing

PRs accepted.

## License

BSD-style © Barret Qin
