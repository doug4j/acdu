package common

import "strings"

//EnsureFowardSlashAtStringEnd makes such "/" is at the end of a string (useful for file processing). Note: this isn't useful for all file systems, but good enough for Unix and Windows-style OSes.
func EnsureFowardSlashAtStringEnd(valDir string) string {
	if !strings.HasSuffix(valDir, "/") {
		valDir = valDir + "/"
	}
	return valDir
}
