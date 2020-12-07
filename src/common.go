package main

//Panics if `err` is not `nil`
//
// @param err Error to be checked for equality to `nil`
func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
