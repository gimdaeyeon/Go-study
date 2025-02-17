package main

func main() {

	names := []string{"홍길동", "이순신", "강감찬"}

	for i, name := range names {
		println(i, name)
	}

}
