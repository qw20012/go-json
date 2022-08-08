package json

import (
	"fmt"
	"testing"

	"github.com/qw20012/go-basic/arr"
)

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

func TestPath(t *testing.T) {
	if container.Path("employees.protected").Data() != "false" {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "false",
			container.Path("employees.protected").Data())
	}
	if container.Path("employees.employee.0.id").String() != "1" {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "1",
			container.Path("employees.employee.0.id").String())
	}
}

func TestSearch(t *testing.T) {
	if container.Search("employees", "address", "country").Data() != "France" {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "France",
			container.Search("employees", "address", "country").Data())
	}
}

func TestExist(t *testing.T) {
	if !container.Exist("employees", "address", "countryCode") {
		t.Fatalf("TestParse expected Type=%s, Got=%v", "true",
			container.Exist("employees", "address", "countryCode"))
	}

	if !container.ExistPath("/employees/address/countryCode") {
		t.Fatalf("TestParse expected Type=%s, Got=%s", "true",
			fmt.Sprintf("%v", container.Exist("employees", "address", "countryCode")))
	}
}
func TestData(t *testing.T) {

	expectedMap := map[string]string{"street": "22 Saint-Lazare", "postalCode": "75003",
		"city": "Paris", "countryCode": "FRA", "country": "France"}
	for key, child := range container.Search("employees", "address").ChildrenMap() {
		//fmt.Printf("Key=>%v, Value=>%v\n", key, child.Data().(string))
		if expectedMap[key] != child.Data().(string) {
			t.Errorf("Child unexpected: %v != %v", expectedMap[key], child.Data().(string))
		}
	}

	expected := []string{"map[first_name:Jeanette id:1 last_name:Penddreth]",
		"map[firstName:Giavani id:2 lastName:Frediani]"}
	// Iterating employee array
	for i, child := range container.Search("employees", "employee").Children() {
		//fmt.Println(child.Data())
		if expected[i] != fmt.Sprintf("%v", child.Data()) {
			t.Errorf("Child unexpected: %v != %v", expected[i], fmt.Sprintf("%v", child.Data()))
		}
	}

	expected = []string{"1", "Jeanette", "Penddreth"}
	// Use index in your search
	for i, child := range container.Search("employees", "employee", "0").Children() {
		if !arr.Contains(expected, fmt.Sprintf("%v", child.Data())) {
			t.Errorf("Contains Child unexpected: %v != %v", expected[i], fmt.Sprintf("%v", child.Data()))
		}
		//fmt.Println(child.Data())
	}
}

func TestGenerateJson(t *testing.T) {
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
}

func TestGenerateArrayJson(t *testing.T) {
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
}

func TestMerge(t *testing.T) {
	jsonParsed1 := Parse(`{"outter":{"value1":"one"}}`)
	jsonParsed2 := Parse(`{"outter":{"inner":{"value3":"three"}},"outter2":{"value2":"two"}}`)

	jsonParsed1.Merge(jsonParsed2)
	expected := "map[outter:map[inner:map[value3:three] value1:one] outter2:map[value2:two]]"
	if fmt.Sprintf("%v", jsonParsed1.Data()) != expected {
		t.Fatalf("TestMerge expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonParsed1.Data()))
	}
}

func TestMergeArray(t *testing.T) {
	jsonParsed1 := Parse(`{"array":["one"]}`)
	jsonParsed2 := Parse(`{"array":["two"]}`)

	jsonParsed1.Merge(jsonParsed2)
	if fmt.Sprintf("%v", jsonParsed1.Data()) != "map[array:[one two]]" {
		t.Fatalf("TestMergeArray expected Type=%s, Got=%v", "map[array:[one two]]",
			fmt.Sprintf("%v", jsonParsed1.Data()))
	}

}

func TestDelete(t *testing.T) {
	jsonParsed := Parse(`{"outter":{"inner":{"value3":"three"}},"outter2":{"value2":"two"}}`)
	jsonParsed.DeletePath("outter.inner.value3")

	expected := "map[outter:map[inner:map[]] outter2:map[value2:two]]"
	if fmt.Sprintf("%v", jsonParsed.Data()) != expected {
		t.Fatalf("TestDelete expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonParsed.Data()))
	}
}

func TestDeleteArray(t *testing.T) {
	jsonParsed := Parse(`{"array":["one","two"]}`)
	jsonParsed.ArrayRemovePath(1, "array")

	expected := "map[array:[one]]"
	if fmt.Sprintf("%v", jsonParsed.Data()) != expected {
		t.Fatalf("TestDeleteArray expected Type=%s, Got=%v", expected,
			fmt.Sprintf("%v", jsonParsed.Data()))
	}
}

func TestContactArray(t *testing.T) {
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
}
