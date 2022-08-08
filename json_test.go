package json

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	jsonStr :=
		`{
    "glossary": {
        "title": "example glossary",
		"GlossDiv": {
            "title": "S",
			"GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
                    },
					"GlossSee": "markup"
                }
            },
            "Nums": 5245243
        }
    }
}`

	container := Parse(jsonStr)
	if container.Path("glossary.title").Data() != "example glossary" {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "example glossary",
			container.Path("glossary.title").Data())
	}

	jsonStr = `
    // Comments
	name:value, 

	"key":"value",
    
	number: 123.5,

	/* Block comments*/
	"another one" : true,`

	js, _ := json.MarshalIndent(Parse(jsonStr).object, "", "    ")

	if !strings.ContainsAny(string(js), "123.5") {
		t.Fatalf("TestParse expected Type=%s, Got=%s", jsonStr, string(js))
	}

}

func TestUnmarshal(t *testing.T) {
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
	//fmt.Println("hello")
	//fmt.Println(aStruct)
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

	aStr := Unmarshal[string]("jsonStr")
	if aStr != "jsonStr" {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", "jsonStr", aStr)
	}
	aInt := Unmarshal[int]("111")
	//fmt.Println(aInt)
	if aInt != 111 {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", "111", aInt)
	}
	aBool := Unmarshal[bool]("true")
	//fmt.Println(aInt)
	if aBool != true {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", "true", aBool)
	}
	aFloat64 := Unmarshal[float64]("1.23")
	//fmt.Println(aInt)
	if aFloat64 != 1.23 {
		t.Fatalf("TestUnmarshal expected Type=%s, Got=%v", "1.23", aFloat64)
	}

	aArray := Unmarshal[[]any](`[101, "202", 303]`)

	if len(aArray) != 3 {
		t.Fatalf("TestUnmarshal expected Type=%v, Got=%v", 3, len(aArray))
	}
	intArray := Unmarshal[[]int](`[101, 202, 303]`)
	if len(intArray) != 3 {
		t.Fatalf("TestUnmarshal expected Type=%v, Got=%v", 3, len(intArray))
	}
	boolArray := Unmarshal[[]bool](`[false, true, false]`)
	if len(boolArray) != 3 {
		t.Fatalf("TestUnmarshal expected Type=%v, Got=%v", 3, len(boolArray))
	}
}

func TestMarshal(t *testing.T) {
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
	//fmt.Println(json)
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
}
