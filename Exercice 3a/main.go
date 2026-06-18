package main

import (
	"errors"
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) DistanceTo(autre Point) float64 {
	return math.Sqrt(math.Pow(autre.X-p.X, 2) + math.Pow(autre.Y-p.Y, 2))
}

func (p Point) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", p.X, p.Y)
}

type Rectangle struct {
	Min Point
	Max Point
}

func NewRectangle(min Point, max Point) (Rectangle, error) {
	if max.X < min.X {
		return Rectangle{}, errors.New("la largeur du rectangle ne peut pas etre negative")
	}

	if max.Y < min.Y {
		return Rectangle{}, errors.New("la hauteur du rectangle ne peut pas etre negative")
	}

	return Rectangle{Min: min, Max: max}, nil
}

func (r Rectangle) Width() float64 {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Height() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) Area() float64 {
	return r.Width() * r.Height()
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width() + r.Height())
}

func (r *Rectangle) Move(dx, dy float64) {
	r.Min.X += dx
	r.Min.Y += dy
	r.Max.X += dx
	r.Max.Y += dy
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle Min=%s, Max=%s, largeur=%.2f, hauteur=%.2f", r.Min, r.Max, r.Width(), r.Height())
}

type Circle struct {
	Center Point
	Radius float64
}

func NewCircle(center Point, radius float64) (Circle, error) {
	if radius < 0 {
		return Circle{}, errors.New("le rayon du cercle ne peut pas etre negatif")
	}

	return Circle{Center: center, Radius: radius}, nil
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

func (c *Circle) Scale(factor float64) error {
	if factor < 0 {
		return errors.New("le facteur d'echelle ne peut pas etre negatif")
	}

	c.Radius *= factor
	return nil
}

func (c Circle) String() string {
	return fmt.Sprintf("Cercle centre=%s, rayon=%.2f", c.Center, c.Radius)
}

func afficherRectangle(rectangle Rectangle) {
	fmt.Println(rectangle)
	fmt.Printf("Largeur: %.2f\n", rectangle.Width())
	fmt.Printf("Hauteur: %.2f\n", rectangle.Height())
	fmt.Printf("Surface: %.2f\n", rectangle.Area())
	fmt.Printf("Perimetre: %.2f\n", rectangle.Perimeter())
}

func afficherCercle(cercle Circle) {
	fmt.Println(cercle)
	fmt.Printf("Surface: %.2f\n", cercle.Area())
	fmt.Printf("Circonference: %.2f\n", cercle.Circumference())
}

func main() {
	fmt.Println("Exercice 1 : Point et Rectangle")
	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 4, Y: 6}
	fmt.Printf("Distance entre %s et %s : %.2f\n", p1, p2, p1.DistanceTo(p2))

	rectangle, err := NewRectangle(Point{X: 0, Y: 0}, Point{X: 5, Y: 3})
	if err != nil {
		fmt.Println("Erreur rectangle :", err)
		return
	}

	afficherRectangle(rectangle)
	rectangle.Move(2, 4)
	fmt.Println("Apres deplacement de dx=2 et dy=4 :")
	fmt.Printf("Min: %s, Max: %s\n", rectangle.Min, rectangle.Max)

	fmt.Println()
	fmt.Println("Exercice 2 : Cercle")
	cercle, err := NewCircle(Point{X: 2, Y: 2}, 4)
	if err != nil {
		fmt.Println("Erreur cercle :", err)
		return
	}

	afficherCercle(cercle)
	if err := cercle.Scale(1.5); err != nil {
		fmt.Println("Erreur mise a l'echelle :", err)
		return
	}

	fmt.Println("Apres mise a l'echelle par 1.5 :")
	afficherCercle(cercle)

	fmt.Println()
	fmt.Println("Exercice 3 : Ameliorations et reflexion")
	if _, err := NewRectangle(Point{X: 5, Y: 0}, Point{X: 1, Y: 3}); err != nil {
		fmt.Println("Validation rectangle invalide :", err)
	}

	if _, err := NewCircle(Point{X: 0, Y: 0}, -2); err != nil {
		fmt.Println("Validation cercle invalide :", err)
	}

	fmt.Println("Approche de validation : utiliser des fonctions constructeur comme NewRectangle et NewCircle qui retournent l'objet et une erreur.")
	fmt.Println("Un receiver de valeur travaille sur une copie de l'objet, ce qui convient aux methodes qui lisent seulement les donnees.")
	fmt.Println("Un receiver de pointeur travaille sur l'objet original, ce qui permet de modifier ses champs.")
	fmt.Println("Move et Scale utilisent donc un receiver pointeur, car elles modifient le rectangle et le cercle.")
	fmt.Println("DistanceTo, Width, Height, Area, Perimeter et Circumference utilisent un receiver valeur, car elles calculent un resultat sans changer l'objet.")
}
