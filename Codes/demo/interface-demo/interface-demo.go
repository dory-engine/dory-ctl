package main

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"golang.org/x/crypto/bcrypt"
	"math"
)

type geometry interface {
	area() float64
	perim() float64
}

type rect struct {
	width, height float64
}
type circle struct {
	radius float64
}

func (r rect) area() float64 {
	return r.width * r.height
}
func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Println(g.perim())
}

func main() {
	r := rect{width: 3, height: 4}
	c := circle{radius: 5}

	measure(r)
	measure(c)

	s := auth.H("j23x6GgfmznBY5QYCbN")
	fmt.Println(auth.RandomKey())
	fmt.Println("###")
	fmt.Println(s)

	bs, _ := bcrypt.GenerateFromPassword([]byte("j23x6GgfmznBY5QYCbN"), 10)
	fmt.Println(string(bs))

	// $2y$10$VD9pK4Egg8TtXn/1IPezXuEou7WLeGnrHl4WmDeBtRxKjH.XGULMG
	// $2a$10$0JaZqjPQEIQ0NUAY7KsUq.M3jhrOLPJzg6MhXGaSzeOOPwX4DRY.W
}
