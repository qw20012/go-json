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
var jsonStr string = `{
	"employees":{
	   "protected":false,
	   "address":{
		  "street":"22 Saint-Lazare",
		  "postalCode":"75003",
		  "city":"Paris",
		  "countryCode":"FRA",
		  "country":"France"
	   },
	   "employee":[
		  {
			 "id":1,
			 "first_name":"Jeanette",
			 "last_name":"Penddreth"
		  },
		  {
			 "id":2,
			 "firstName":"Giavani",
			 "lastName":"Frediani"
		  }
	   ]
	}
 }`

var container *Container = Parse(jsonStr)
```
### container.Path

Path searches the wrapped structure following a path in dot or forward slash notation, segments of this path are searched according to the same rules as Search.
```
	if container.Path("employees.protected").Data() != "false" {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "false",
			container.Path("employees.protected").Data())
	}
	
	if container.Path("employees.employee.0.id").String() != "1" {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "1",
			container.Path("employees.employee.0.id").String())
	}
```
### container.Search

Search attempts to find and return an object within the wrapped structure by following a provided hierarchy of field names to locate the target.
```
	if container.Search("employees", "address", "country").Data() != "France" {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "France",
			container.Search("employees", "address", "country").Data())
	}
```
### container.Exists

Exists checks whether a field exists within the hierarchy.
```
	if !container.Exists("employees", "address", "countryCode") {
		t.Fatalf("TestParse expected Type=%s, Got=%v", "true",
			container.Exists("employees", "address", "countryCode"))
	}
```
### container.ExistPath

ExistPath checks whether a dot or forward slash notation path exists.
```
	if !container.ExistPath("/employees/address/countryCode") {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "true",
			fmt.Sprintf("%v", container.Exists("employees", "address", "countryCode")))
	}
```
### container.ChildrenMap

ChildrenMap returns a map of all the children of an object element. IF the underlying value isn't a object then an empty map is returned.
```
	expectedMap := map[string]string{"street": "22 Saint-Lazare", "postalCode": "75003",
		"city": "Paris", "countryCode": "FRA", "country": "France"}
	for key, child := range container.Search("employees", "address").ChildrenMap() {
		//fmt.Printf("Key=>%v, Value=>%v\n", key, child.Data().(string))
		if expectedMap[key] != child.Data().(string) {
			t.Errorf("Child unexpected: %v != %v", expectedMap[key], child.Data().(string))
		}
	}
```
### container.Children

Children returns a slice of all children of an array element. This also works for objects, however, the children returned for an object will be in a random order and you lose the names of the returned objects this way. If the underlying container value isn't an array or map nil is returned.
```
	expected := []string{"map[first_name:Jeanette id:1 last_name:Penddreth]",
		"map[firstName:Giavani id:2 lastName:Frediani]"}
	// Iterating employee array
	for i, child := range container.Search("employees", "employee").Children() {
		//fmt.Println(child.Data())
		if expected[i] != fmt.Sprintf("%v", child.Data()) {
			t.Errorf("Child unexpected: %v != %v", expected[i], fmt.Sprintf("%v", child.Data()))
		}
	}
```
### GenerateJson

```
	jsonObj := New()
	// or gabs.Wrap(jsonObject) to work on an existing map[string]interface{}

	jsonObj.Set(10, "outter", "inner", "value")
	jsonObj.SetPath(20, "outter.inner.value2")
	jsonObj.Set(30, "outter", "inner2", "value3")

	expected := "map[outter:map[inner:map[value:10 value2:20] inner2:map[value3:30]]]"
	if fmt.Sprintf("%v", jsonObj.Data()) != expected {
		t.Fatalf("TestGenerateJson expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonObj.Data()))
	}
```
### GenerateArrayJson

```
	jsonObj := New()

	jsonObj.Array("foo", "array")
	// Or .ArrayP("foo.array")

	jsonObj.ArrayAppend(10, "foo", "array")
	jsonObj.ArrayAppend(20, "foo", "array")
	jsonObj.ArrayAppend(30, "foo", "array")

	expected := "map[foo:map[array:[10 20 30]]]"
	if fmt.Sprintf("%v", jsonObj.Data()) != expected {
		t.Fatalf("TestGenerateArrayJson expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonObj.Data()))
	}
```
### container.Merge

Merge a source object into an existing destination object. When a collision is found within the merged structures (both a source and destination object contain the same non-object keys) the result will be an array containing both values, where values that are already arrays will be expanded into the resulting array.
```
	jsonParsed1 := Parse(`{"outter":{"value1":"one"}}`)
	jsonParsed2 := Parse(`{"outter":{"inner":{"value3":"three"}},"outter2":{"value2":"two"}}`)

	jsonParsed1.Merge(jsonParsed2)
	expected := "map[outter:map[inner:map[value3:three] value1:one] outter2:map[value2:two]]"
	if fmt.Sprintf("%v", jsonParsed1.Data()) != expected {
		t.Fatalf("TestMerge expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonParsed1.Data()))
	}
```
### container.DeletePath

DeletePath deletes an element at a path using dot or forward slash notation, an error is returned if the element does not exist.
```
	jsonParsed := Parse(`{"outter":{"inner":{"value3":"three"}},"outter2":{"value2":"two"}}`)
	jsonParsed.DeletePath("outter.inner.value3")

	expected := "map[outter:map[inner:map[]] outter2:map[value2:two]]"
	if fmt.Sprintf("%v", jsonParsed.Data()) != expected {
		t.Fatalf("TestDelete expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonParsed.Data()))
	}
```
### container.ArrayRemovePath

ArrayRemoveP attempts to remove an element identified by an index from a JSON array at a path using dot or forward slash notation.
```
	jsonParsed := Parse(`{"array":["one","two"]}`)
	jsonParsed.ArrayRemovePath(1, "array")

	expected := "map[array:[one]]"
	if fmt.Sprintf("%v", jsonParsed.Data()) != expected {
		t.Fatalf("TestDeleteArray expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonParsed.Data()))
	}
```
### container.ArrayConcat

ArrayConcat attempts to append a value onto a JSON array at a path. If the target is not a JSON array then it will be converted into one, with its original contents set to the first element of the array.
```
	jsonObj := New()

	jsonObj.Array("foo", "array")
	// Or .ArrayP("foo.array")

	jsonObj.ArrayConcat(10, "foo", "array")
	jsonObj.ArrayConcat([]interface{}{20, 30}, "foo", "array")

	expected := "map[foo:map[array:[10 20 30]]]"
	if fmt.Sprintf("%v", jsonObj.Data()) != expected {
		t.Fatalf("TestGenerateArrayJson expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonObj.Data()))
	}
```

## Contributing

PRs accepted.

## License

BSD-style © Barret Qin
