# go-json

- Json for humans. It is very suitable for the config json file.
	+ Single and multi-line comments are allowed.
	+ Trailing comma is allowed.
	+ Quotes can be ignored when string contains no space.
	+ Outer brace can be ignored.
- Helpful wrapper for navigating hierarchies of map[string]any objects.
- Encoding and decoding of JSON like Package json.

## Usage

### json.ParseFile

Parse json string from given file. This can be used to process configuraton file.
```
	config, err := ParseFile[map[string]any]("config.json")
	if err != nil {
		t.Fatalf("TestParseFile failed")
	}

	if config["key"].(string) != "value" {
		t.Fatalf("TestParseFile failed Type=false, Got=true")
	}
```
The content in file config.json as below.
```
	// '{' Outer barce can be ignored.
	
	// Quotes can be ignored when string contains no space.
	name:value, 

	"key":"value",
    
	number: 123.5,

	/* 
	  Block comments
	  Trailing comma is allowed.
	*/
	"another one" : true,
	// '}' Outer barce can be ignored.
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
