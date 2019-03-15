package utility

func SepArg(sep string) string {
	if sep == "t" || sep == "\\t" || sep == "s" || sep == "\\s" {
		sep = "\t"
	}
	return sep
}
