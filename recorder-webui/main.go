package main

func main() {

	test := &PerfTest{state: &State{}}
	test.Start()
    test.Await()
}
