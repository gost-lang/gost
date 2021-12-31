package object

import "fmt"

const LIBRARY_MODULE = "LIBRARY_MODULE"

// LibraryModule objects consist of a slice of LibraryFunctions.
type LibraryModule struct {
	Name    string
	Methods map[string]*LibraryFunction
}

func (object *LibraryModule) Accept(v Visitor) {
	v.visitLibraryModule(object)
}

// String represents the library module's value as a string.
func (libraryModule *LibraryModule) String() string {
	return fmt.Sprintf("library module {%s}", libraryModule.Name)
}

// Type returns the library module object type.
func (libraryModule *LibraryModule) Type() Type {
	return LIBRARY_MODULE
}
