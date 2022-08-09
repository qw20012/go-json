# go-json

- Json for humans. It is very suitable for the config json file.
	+ Single and multi-line comments are allowed.
	+ Trailing comma is allowed.
	+ Quotes can be ignored when string contains no space.
	+ Outer brace can be ignored.
- Encoding and decoding of JSON like Package json.
- Helpful wrapper for navigating hierarchies of map[string]any objects.

## Usage

### json.UnmarshalFile

Unmarshal json string from given file. This can be used to process configuraton file.
```
	config, err := UnmarshalFile[map[string]any]("config.json")
	if err != nil {
		t.Fatalf("UnmarshalFile failed")
	}

	if config["key"].(string) != "value" {
		t.Fatalf("UnmarshalFile failed Type=false, Got=true")
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
### json.Unmarshal

Unmarshal json string into given generic type (T).
```
	type Book struct {
		Name  *string
		Pages *int
		Arr   *[]string
	}
	
	type Library struct {
		Ptr     *[]*int
		Book    Book //Books string
		BookPtr *Book
		Books   []Book
		Count   int
		IsNew   bool
		Price   float32
	}
	
	jsonStr := `
	    // Comments
		Ptr:[1,2,3],Book:{Name:"red", Pages:100, Arr:["1","2","a"]},BookPtr:{Name:"red", Pages:100, Arr:["1","2","a"]},Books:[{Name:"red", Pages:100, Arr:["1","2","a"]},{Name:"red", Pages:100, Arr:["1","2","a"]}], Count:3, IsNew:true, Price:1.234`
	
	aStruct := Unmarshal[Library](jsonStr)
	
	if *aStruct.Book.Name != "red" {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", "red",
			fmt.Sprintf("%v", aStruct))
	}

	if *aStruct.BookPtr.Pages != 100 {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", "100",
			*aStruct.BookPtr.Pages)
	}

	if (*aStruct.BookPtr.Arr)[0] != "1" {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", "1",
			(*aStruct.BookPtr.Arr)[0])
	}
	
	if *((*aStruct.Ptr)[0]) != 1 {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", "1",
			fmt.Sprintf("%v", aStruct))
	}
	
	jsonStr = `
	// Comments
	Books:{Name:"red", Pages:100}, Count:3, IsNew:true, Price:1.234`
	aMap := Unmarshal[map[string]any](jsonStr)
	expected := "map[Books:map[Name:red Pages:100] Count:3 IsNew:true Price:1.234]"
	if fmt.Sprintf("%v", aMap) != expected {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", aMap))
	}

	jsonStr = `
	    /* Comments */
		Books:1, Count:2, IsNew:3`
	bMap := Unmarshal[map[string]int](jsonStr)
	expected = "map[Books:1 Count:2 IsNew:3]"
	if fmt.Sprintf("%v", bMap) != expected {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", bMap))
	}
```
### json.Marshal

Marshal given struct/map or their pointer to json string.
```
	type Book struct {
		Name  *string
		Pages *int `json:"pages"`
		Arr   *[]*string
	}
	
	type Library struct {
		Ptr     *[]*int
		Book    Book // Books string
		BookPtr *Book
		Books   []Book
		Count   int
		IsNew   bool
		Price   float32
	}
	
	aInt := 1
	aIntArr := []*int{&aInt}

	bookName := "book name"
	pages := 10
	arr := []*string{&bookName}

	var book = Book{Name: &bookName, Pages: &pages, Arr: &arr}
	var lib = Library{Ptr: &aIntArr, Book: book, BookPtr: &book, Books: []Book{book, book}, Count: 1, IsNew: true, Price: 1.2}
	
	json := Marshal(lib)
	if !strings.Contains(json, "1.2") {
		t.Fatalf("TestUnmarshal expected Type=%v, Got=%v", "1.2", json)
	}

	libPtr := &Library{Ptr: &aIntArr, Book: book, BookPtr: &book, Books: []Book{book, book}, Count: 1, IsNew: true, Price: 1.1}
	
	json = Marshal(libPtr)
	if !strings.Contains(json, "1.1") {
		t.Fatalf("TestUnmarshal expected Type=%v, Got=%v", "1.1", json)
	}
	
	aMap := map[string]any{}
	aMap["int"] = 10
	aMap["string"] = "this is a string"
	aMap["bool"] = false
	aMap["bMap"] = map[string]any{"a": 1.2}
	
	json = Marshal(aMap)
	if !strings.Contains(json, "1.2") {
		t.Fatalf("TestUnmarshal expected Type=%v, Got=%v", "1.2", json)
	}
```
### Container

Helpful wrapper for navigating hierarchies of map[string]any objects.
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
