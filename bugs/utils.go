package bugs

func check(e error) {
	if e != nil {
	//	fmt.Fprintf(os.Stderr, "err: %s", err.Error())
	//	return NoConfigError
		panic(e)
	}
}
