package json

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrOutOfBounds indicates an index was out of bounds.
	ErrOutOfBounds = errors.New("out of bounds")

	// ErrNotObjOrArray is returned when a target is not an object or array type
	// but needs to be for the intended operation.
	ErrNotObjOrArray = errors.New("not an object or array")

	// ErrNotObj is returned when a target is not an object but needs to be for
	// the intended operation.
	ErrNotObj = errors.New("not an object")

	// ErrInvalidQuery is returned when a seach query was not valid.
	ErrInvalidQuery = errors.New("invalid search query")

	// ErrNotArray is returned when a target is not an array but needs to be for
	// the intended operation.
	ErrNotArray = errors.New("not an array")

	// ErrPathCollision is returned when creating a path failed because an
	// element collided with an existing value.
	ErrPathCollision = errors.New("encountered value collision whilst building path")

	// ErrNotFound is returned when a query leaf is not found.
	ErrNotFound = errors.New("field not found")
)

// Container references a specific element within a wrapped structure. See to gabs.
type Container struct {
	object any
}

func (g *Container) searchStrict(allowWildcard bool, hierarchy ...string) (*Container, error) {
	object := g.Data()

	for target := 0; target < len(hierarchy); target++ {
		pathSeg := hierarchy[target]
		if mmap, ok := object.(map[string]any); ok {
			object, ok = mmap[pathSeg]
			if !ok {
				return nil, fmt.Errorf("failed to resolve path segment '%v': key '%v' was not found", target, pathSeg)
			}

		} else if marray, ok := object.([]any); ok {
			if allowWildcard && pathSeg == "*" {
				tmpArray := []any{}
				for _, val := range marray {
					if (target + 1) >= len(hierarchy) {
						tmpArray = append(tmpArray, val)
					} else if res := Wrap(val).Search(hierarchy[target+1:]...); res != nil {
						tmpArray = append(tmpArray, res.Data())
					}
				}
				if len(tmpArray) == 0 {
					return nil, nil
				}
				return &Container{tmpArray}, nil
			}
			index, err := strconv.Atoi(pathSeg)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve path segment '%v': found array but segment value '%v' could not be parsed into array index: %v", target, pathSeg, err)
			}
			if index < 0 {
				return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' is invalid", target, pathSeg)
			}
			if len(marray) <= index {
				return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' exceeded target array size of '%v'", target, pathSeg, len(marray))
			}
			object = marray[index]
		} else {
			return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was not found", target, pathSeg)
		}
	}

	return &Container{object}, nil
}

// Search attempts to find and return an object within the wrapped structure by
// following a provided hierarchy of field names to locate the target.
//
// If the search encounters an array then the next hierarchy field name must be
// either a an integer which is interpreted as the index of the target, or the
// character '*', in which case all elements are searched with the remaining
// search hierarchy and the results returned within an array.
func (g *Container) Search(hierarchy ...string) *Container {
	c, _ := g.searchStrict(true, hierarchy...)
	return c
}

// Path searches the wrapped structure following a path in dot or forward slash notation,
// segments of this path are searched according to the same rules as Search.
//
// Because the characters '~' (%x7E), '.' (%x2E) and '/' have special meanings in paths,
// '~' needs to be encoded as '~0' and '.' or '/' needs to be encoded as
// '~1' when these characters appear in a reference key.
func (g *Container) Path(path string) *Container {
	return g.Search(PathToSlice(path)...)
}

// Exists checks whether a field exists within the hierarchy.
func (g *Container) Exist(hierarchy ...string) bool {
	return g.Search(hierarchy...) != nil
}

// ExistPath checks whether a dot or forward slash notation path exists.
func (g *Container) ExistPath(path string) bool {
	return g.Exist(PathToSlice(path)...)
}

// Children returns a slice of all children of an array element. This also works
// for objects, however, the children returned for an object will be in a random
// order and you lose the names of the returned objects this way. If the
// underlying container value isn't an array or map nil is returned.
func (g *Container) Children() []*Container {
	if array, ok := g.Data().([]any); ok {
		children := make([]*Container, len(array))
		for i := 0; i < len(array); i++ {
			children[i] = &Container{array[i]}
		}
		return children
	}
	if mmap, ok := g.Data().(map[string]any); ok {
		children := []*Container{}
		for _, obj := range mmap {
			children = append(children, &Container{obj})
		}
		return children
	}
	return nil
}

// ChildrenMap returns a map of all the children of an object element. IF the
// underlying value isn't a object then an empty map is returned.
func (g *Container) ChildrenMap() map[string]*Container {
	if mmap, ok := g.Data().(map[string]any); ok {
		children := make(map[string]*Container, len(mmap))
		for name, obj := range mmap {
			children[name] = &Container{obj}
		}
		return children
	}
	return map[string]*Container{}
}

// New creates a new Container JSON object.
func New() *Container {
	return &Container{map[string]any{}}
}

// Wrap an already unmarshalled JSON object (or a new map[string]any)
// into a *Container.
func Wrap(root any) *Container {
	return &Container{root}
}

// Data returns the underlying value of the target element in the wrapped
// structure.
func (g *Container) Data() any {
	if g == nil {
		return nil
	}
	return g.object
}

// String marshals an element to a JSON formatted string.
func (g *Container) String() string {
	return g.Data().(string)
}

// PathToSlice returns a slice of path segments parsed out of a dot or forward slash path.
//
// Because '.' (%x2E) or '/' is the segment separator, it must be encoded as '~1'
// if it appears in the reference key. Likewise, '~' (%x7E) must be encoded
// as '~0' since it is the escape character for encoding '.' or '/'.
func PathToSlice(path string) []string {
	if len(path) == 0 {
		return nil
	}
	if path == "/" || path == "." {
		return []string{""}
	}

	var hierarchy []string
	if path[0] != '/' {
		hierarchy = strings.Split(path, ".")
		for i, v := range hierarchy {
			v = strings.Replace(v, "~1", ".", -1)
			v = strings.Replace(v, "~0", "~", -1)
			hierarchy[i] = v
		}
	} else {
		hierarchy := strings.Split(path, "/")[1:]
		for i, v := range hierarchy {
			v = strings.Replace(v, "~1", "/", -1)
			v = strings.Replace(v, "~0", "~", -1)
			hierarchy[i] = v
		}
	}

	return hierarchy
}

//------------------------------------------------------------------------------

// Set attempts to set the value of a field located by a hierarchy of field
// names. If the search encounters an array then the next hierarchy field name
// is interpreted as an integer index of an existing element, or the character
// '-', which indicates a new element appended to the end of the array.
//
// Any parts of the hierarchy that do not exist will be constructed as objects.
// This includes parts that could be interpreted as array indexes.
//
// Returns a container of the new value or an error.
func (g *Container) Set(value any, hierarchy ...string) (*Container, error) {
	if g == nil {
		return nil, errors.New("failed to resolve path, container is nil")
	}
	if len(hierarchy) == 0 {
		g.object = value
		return g, nil
	}
	if g.object == nil {
		g.object = map[string]any{}
	}
	object := g.object

	for target := 0; target < len(hierarchy); target++ {
		pathSeg := hierarchy[target]
		if mmap, ok := object.(map[string]any); ok {
			if target == len(hierarchy)-1 {
				object = value
				mmap[pathSeg] = object
			} else if object = mmap[pathSeg]; object == nil {
				mmap[pathSeg] = map[string]any{}
				object = mmap[pathSeg]
			}
		} else if marray, ok := object.([]any); ok {
			if pathSeg == "-" {
				if target < 1 {
					return nil, errors.New("unable to append new array index at root of path")
				}
				if target == len(hierarchy)-1 {
					object = value
				} else {
					object = map[string]any{}
				}
				marray = append(marray, object)
				if _, err := g.Set(marray, hierarchy[:target]...); err != nil {
					return nil, err
				}
			} else {
				index, err := strconv.Atoi(pathSeg)
				if err != nil {
					return nil, fmt.Errorf("failed to resolve path segment '%v': found array but segment value '%v' could not be parsed into array index: %v", target, pathSeg, err)
				}
				if index < 0 {
					return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' is invalid", target, pathSeg)
				}
				if len(marray) <= index {
					return nil, fmt.Errorf("failed to resolve path segment '%v': found array but index '%v' exceeded target array size of '%v'", target, pathSeg, len(marray))
				}
				if target == len(hierarchy)-1 {
					object = value
					marray[index] = object
				} else if object = marray[index]; object == nil {
					return nil, fmt.Errorf("failed to resolve path segment '%v': field '%v' was not found", target, pathSeg)
				}
			}
		} else {
			return nil, ErrPathCollision
		}
	}
	return &Container{object}, nil
}

// SetPath sets the value of a field at a path using dot or forward slash notation, any parts
// of the path that do not exist will be constructed, and if a collision occurs
// with a non object type whilst iterating the path an error is returned.
func (g *Container) SetPath(value any, path string) (*Container, error) {
	return g.Set(value, PathToSlice(path)...)
}

//------------------------------------------------------------------------------

/*
Array modification/search - Keeping these options simple right now, no need for
anything more complicated since you can just cast to []any, modify and
then reassign with Set.
*/

// ArrayAppend attempts to append a value onto a JSON array at a path. If the
// target is not a JSON array then it will be converted into one, with its
// original contents set to the first element of the array.
func (g *Container) ArrayAppend(value any, hierarchy ...string) error {
	if array, ok := g.Search(hierarchy...).Data().([]any); ok {
		array = append(array, value)
		_, err := g.Set(array, hierarchy...)
		return err
	}

	newArray := []any{}
	if d := g.Search(hierarchy...).Data(); d != nil {
		newArray = append(newArray, d)
	}
	newArray = append(newArray, value)

	_, err := g.Set(newArray, hierarchy...)
	return err
}

// ArrayAppendPath attempts to append a value onto a JSON array at a path using dot
// notation. If the target is not a JSON array then it will be converted into
// one, with its original contents set to the first element of the array.
func (g *Container) ArrayAppendPath(value any, path string) error {
	return g.ArrayAppend(value, PathToSlice(path)...)
}

// Array creates a new JSON array at a path. Returns an error if the path
// contains a collision with a non object type.
func (g *Container) Array(hierarchy ...string) (*Container, error) {
	return g.Set([]any{}, hierarchy...)
}

// ArrayP creates a new JSON array at a path using dot or forward slash notation. Returns an
// error if the path contains a collision with a non object type.
func (g *Container) ArrayPath(path string) (*Container, error) {
	return g.Array(PathToSlice(path)...)
}

// Delete an element at a path, an error is returned if the element does not
// exist or is not an object. In order to remove an array element please use
// ArrayRemove.
func (g *Container) Delete(hierarchy ...string) error {
	if g == nil || g.object == nil {
		return ErrNotObj
	}
	if len(hierarchy) == 0 {
		return ErrInvalidQuery
	}

	object := g.object
	target := hierarchy[len(hierarchy)-1]
	if len(hierarchy) > 1 {
		object = g.Search(hierarchy[:len(hierarchy)-1]...).Data()
	}

	if obj, ok := object.(map[string]any); ok {
		if _, ok = obj[target]; !ok {
			return ErrNotFound
		}
		delete(obj, target)
		return nil
	}
	if array, ok := object.([]any); ok {
		if len(hierarchy) < 2 {
			return errors.New("unable to delete array index at root of path")
		}
		index, err := strconv.Atoi(target)
		if err != nil {
			return fmt.Errorf("failed to parse array index '%v': %v", target, err)
		}
		if index >= len(array) {
			return ErrOutOfBounds
		}
		if index < 0 {
			return ErrOutOfBounds
		}
		array = append(array[:index], array[index+1:]...)
		g.Set(array, hierarchy[:len(hierarchy)-1]...)
		return nil
	}
	return ErrNotObjOrArray
}

// DeleteP deletes an element at a path using dot or forward slash notation, an error is returned
// if the element does not exist.
func (g *Container) DeletePath(path string) error {
	return g.Delete(PathToSlice(path)...)
}

// MergeFn merges two objects using a provided function to resolve collisions.
//
// The collision function receives two any arguments, destination (the
// original object) and source (the object being merged into the destination).
// Which ever value is returned becomes the new value in the destination object
// at the location of the collision.
func (g *Container) MergeFn(source *Container, collisionFn func(destination, source any) any) error {
	var recursiveFnc func(map[string]any, []string) error
	recursiveFnc = func(mmap map[string]any, path []string) error {
		for key, value := range mmap {
			newPath := append(path, key)
			if g.Exist(newPath...) {
				existingData := g.Search(newPath...).Data()
				switch t := value.(type) {
				case map[string]any:
					switch existingVal := existingData.(type) {
					case map[string]any:
						if err := recursiveFnc(t, newPath); err != nil {
							return err
						}
					default:
						if _, err := g.Set(collisionFn(existingVal, t), newPath...); err != nil {
							return err
						}
					}
				default:
					if _, err := g.Set(collisionFn(existingData, t), newPath...); err != nil {
						return err
					}
				}
			} else {
				// path doesn't exist. So set the value
				if _, err := g.Set(value, newPath...); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if mmap, ok := source.Data().(map[string]any); ok {
		return recursiveFnc(mmap, []string{})
	}
	return nil
}

// Merge a source object into an existing destination object. When a collision
// is found within the merged structures (both a source and destination object
// contain the same non-object keys) the result will be an array containing both
// values, where values that are already arrays will be expanded into the
// resulting array.
//
// It is possible to merge structures will different collision behaviours with
// MergeFn.
func (g *Container) Merge(source *Container) error {
	return g.MergeFn(source, func(dest, source any) any {
		destArr, destIsArray := dest.([]any)
		sourceArr, sourceIsArray := source.([]any)
		if destIsArray {
			if sourceIsArray {
				return append(destArr, sourceArr...)
			}
			return append(destArr, source)
		}
		if sourceIsArray {
			return append(append([]any{}, dest), sourceArr...)
		}
		return []any{dest, source}
	})
}

// ArrayConcat attempts to append a value onto a JSON array at a path. If the
// target is not a JSON array then it will be converted into one, with its
// original contents set to the first element of the array.
//
// ArrayConcat differs from ArrayAppend in that it will expand a value type
// []any during the append operation, resulting in concatenation of each
// element, rather than append as a single element of []any.
func (g *Container) ArrayConcat(value any, hierarchy ...string) error {
	var array []any
	if d := g.Search(hierarchy...).Data(); d != nil {
		if targetArray, ok := d.([]any); !ok {
			// If the data exists, and it is not a slice of interface,
			// append it as the first element of our new array.
			array = append(array, d)
		} else {
			// If the data exists, and it is a slice of interface,
			// assign it to our variable.
			array = targetArray
		}
	}

	switch v := value.(type) {
	case []any:
		// If we have been given a slice of interface, expand it when appending.
		array = append(array, v...)
	default:
		array = append(array, v)
	}

	_, err := g.Set(array, hierarchy...)

	return err
}

// ArrayConcatPath attempts to append a value onto a JSON array at a path using dot
// notation. If the target is not a JSON array then it will be converted into one,
// with its original contents set to the first element of the array.
//
// ArrayConcatPath differs from ArrayAppendPath in that it will expand a value type
// []any during the append operation, resulting in concatenation of each
// element, rather than append as a single element of []any.
func (g *Container) ArrayConcatPath(value any, path string) error {
	return g.ArrayConcat(value, PathToSlice(path)...)
}

// ArrayRemove attempts to remove an element identified by an index from a JSON
// array at a path.
func (g *Container) ArrayRemove(index int, hierarchy ...string) error {
	if index < 0 {
		return ErrOutOfBounds
	}
	array, ok := g.Search(hierarchy...).Data().([]any)
	if !ok {
		return ErrNotArray
	}
	if index < len(array) {
		array = append(array[:index], array[index+1:]...)
	} else {
		return ErrOutOfBounds
	}
	_, err := g.Set(array, hierarchy...)
	return err
}

// ArrayRemoveP attempts to remove an element identified by an index from a JSON
// array at a path using dot or forward slash notation.
func (g *Container) ArrayRemovePath(index int, path string) error {
	return g.ArrayRemove(index, PathToSlice(path)...)
}
