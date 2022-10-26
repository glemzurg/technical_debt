package technical_debt

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

type packageFolder struct {
	name       string        // Name of this package.
	importPath string        // Package relative to src.
	files      []packageFile // The files in this package.
}

func (f packageFolder) String() (output string) {

	var fileStrings []string
	for _, file := range f.files {
		fileStrings = append(fileStrings, file.String())
	}

	return fmt.Sprintf("%s (%s)\n\n%s\n", f.name, f.importPath, strings.Join(fileStrings, "\n"))
}

type packageFile struct {
	name         string // Name of this file.
	imports      map[string]fileImport
	declarations []fileDeclaration
	unresolved   []fileUnresolved
}

func (f packageFile) String() (output string) {

	var importString string
	if len(f.imports) > 0 {
		var importStrings []string
		for _, imp := range f.imports {
			importStrings = append(importStrings, imp.String())
		}
		importString = fmt.Sprintf("IMPORTS:\n%s\n", strings.Join(importStrings, "\n"))
	}

	var declarationString string
	if len(f.declarations) > 0 {
		var declarationStrings []string
		for _, declartion := range f.declarations {
			declarationStrings = append(declarationStrings, declartion.String())
		}
		declarationString = fmt.Sprintf("DECLARATIONS:\n%s\n", strings.Join(declarationStrings, "\n"))
	}

	var unresolvedString string
	if len(f.unresolved) > 0 {
		var unresolvedStrings []string
		for _, unresolved := range f.unresolved {
			unresolvedStrings = append(unresolvedStrings, unresolved.String())
		}
		unresolvedString = fmt.Sprintf("UNRESOLVED:\n%s\n", strings.Join(unresolvedStrings, "\n"))
	}

	return fmt.Sprintf("%s\n%s%s%s", f.name, importString, declarationString, unresolvedString)
}

type fileImport struct {
	name      string
	path      string
	inProject bool
}

func (i fileImport) String() (output string) {
	return fmt.Sprintf("\t%s %s\n", i.name, i.path)
}

type fileDeclaration struct {
	kind string
	name string
}

func (d fileDeclaration) String() (output string) {
	return fmt.Sprintf("\t%s (%s)\n", d.name, d.kind)
}

type fileUnresolved struct {
	packageName string // Relevant package.
	name        string // The active item if in this package or a project package. "*" if for an out-of-project package.
}

func (u fileUnresolved) String() (output string) {
	return fmt.Sprintf("\t%s.%s\n", u.packageName, u.name)
}

type unresolvedName struct {
	offset int    // The byte offset into the source file with the reference.
	name   string // The text string of the unresolved refernce. May only be the package name of a longer identifier.
}

// ProcessPackage processes all the tokens of a single package.
func ProcessPackage(gopath string, packagePaths, projectPaths []string, includeTests bool) (folders []packageFolder, err error) {
	var ok bool

	// Go through every package.
	for _, packagePath := range packagePaths {
		fset := token.NewFileSet() // positions are relative to fset

		// Parse the package.
		p, err := parser.ParseDir(fset, gopath+"/"+packagePath, nil, 0)
		if err != nil {
			return nil, Error(err)
		}

		for _, parsedPackage := range p {
			var folder packageFolder
			folder.name = parsedPackage.Name
			folder.importPath = packagePath

			for filename, parsedFile := range parsedPackage.Files {
				var file packageFile

				// Continue with this file if we are including both tests and code files, or this is not a test.
				if includeTests || !strings.HasSuffix(filename, "_test.go") {
					file.name = filepath.Base(filename)

					// We'll need the text of the file later.
					var data []byte
					if data, err = ioutil.ReadFile(filename); err != nil {
						return nil, Error(err)
					}
					var fileText string = string(data)

					// Get all the imports.
					file.imports = map[string]fileImport{}
					for _, s := range parsedFile.Imports {

						// Get the import path in a usuable format.
						importPath, err := strconv.Unquote(s.Path.Value)
						if err != nil {
							return nil, Error(err)
						}

						// Only continue with import paths that are in our projects.
						inProject := false
						for _, projectPath := range projectPaths {
							if strings.HasPrefix(importPath, projectPath) {
								inProject = true
								break
							}
						}

						// If this is in our projects we want to look for dependencies.
						if inProject {
							var importName string = filepath.Base(importPath)
							if s.Name != nil && s.Name.Name != "." {
								importName = s.Name.Name
							}
							file.imports[importName] = fileImport{name: importName, path: importPath, inProject: true}
						}
					}

					for _, object := range parsedFile.Scope.Objects {
						file.declarations = append(file.declarations, fileDeclaration{kind: object.Kind.String(), name: object.Name})
					}

					// All unresolved.
					var rawUnresolvedNames []unresolvedName
					for _, unresolved := range parsedFile.Unresolved {
						rawUnresolvedNames = append(rawUnresolvedNames, unresolvedName{
							offset: fset.Position(unresolved.NamePos).Offset,
							name:   unresolved.Name,
						})
					}

					// Drop anything that is definitely not a connection to another code file.
					// Imports to non-project packages can disovered just with the import data structure.
					var validUnresolvedNames []unresolvedName
					for _, unresolved := range rawUnresolvedNames {
						switch unresolved.name {
						case "error", "nil", "bool", "string", "true", "false", "uint", "int", "uint64", "int64", "float64":
							// Not a name to keep, not a link to another file.
						default:
							validUnresolvedNames = append(validUnresolvedNames, unresolved)
						}
					}

					// Collapse to an actually distinct set of accurate references.
					var unresolvedSet map[string]fileUnresolved = map[string]fileUnresolved{}
					for _, unresolved := range validUnresolvedNames {
						// Is this unresolved a package?
						var theImport fileImport
						if theImport, ok = file.imports[unresolved.name]; ok {
							if theImport.inProject {
								// There is a reference to an in-project package.

								// This will be a reference like "package.ExportedMember"
								// At the moment we only have "package", gather true name that is unresolved.
								var packageName string = unresolved.name

								// There should be a period.
								var periodOffset int = unresolved.offset + len(unresolved.name)
								if rune(fileText[periodOffset]) != '.' {
									panic(`expected period when looking for unresolved reference`)
								}

								// Build out the rest of the name.
								var unresolvedNameRunes []rune
								var keepReading bool = true
								for i := periodOffset + 1; keepReading; i++ {
									var identifierRune rune
									if identifierRune, keepReading = validIdentifierRune(fileText[i]); keepReading {
										unresolvedNameRunes = append(unresolvedNameRunes, identifierRune)
									}
								}
								var unresolvedName string = string(unresolvedNameRunes)

								// If unresolved is ever empty we have an unexpected case.
								if unresolvedName == "" {
									panic(`unexpected empty name`)
								}

								// Add what we have.
								unresolvedSet[packageName+"."+unresolvedName] = fileUnresolved{packageName: theImport.path, name: unresolvedName}

							} else {
								// Not in the project package.
								// We know this is in the imports so no need to record it here.
							}
						} else {
							// Not an import. This is part of this package.
							unresolvedSet[folder.name+"."+unresolved.name] = fileUnresolved{packageName: folder.importPath, name: unresolved.name}
						}
					}

					// Put them into the proper format.
					for _, unresolved := range unresolvedSet {
						file.unresolved = append(file.unresolved, unresolved)
					}

					folder.files = append(folder.files, file)
				}
			}

			folders = append(folders, folder)
		}
	}

	return folders, nil
}

func validIdentifierRune(b byte) (r rune, keepReading bool) {
	switch rune(b) {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		return rune(b), true
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		return rune(b), true
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return rune(b), true
	case '_':
		return rune(b), true
	}
	return rune(b), false
}
