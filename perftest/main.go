package main

func main() {

	test := &PerfTest{state: &State{}}
	test.RegisterNumber()

	test.Subscribe()
    test.Await()
}
