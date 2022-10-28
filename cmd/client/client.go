package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dvher/Tarea2SD/pkg/coordinates"
	"github.com/dvher/Tarea2SD/pkg/miembro"
	"github.com/dvher/Tarea2SD/pkg/venta"
)

var (
	currMember string
)

func generateCoords(min, max float64) [2]float64 {

	var coords [2]float64

	for i := range coords {
		coords[i] = min + rand.Float64()*(max-min)
	}

	return coords

}

func sendLocation() {
	ticker := time.NewTicker(time.Duration(rand.Intn(70)) * time.Second)

	for {

		if currMember == "" {
			continue
		}

		xy := generateCoords(0, 100)
		c := coordinates.Coordinates{
			Coords:  xy,
			Miembro: currMember,
		}
		sendCoords(c)
		<-ticker.C
	}
}

func sendVentas(v venta.Venta) {

	txt, err := v.JSON()
	if err != nil {
		log.Panic(err)
	}

	body := strings.NewReader(string(txt))

	if _, err = http.Post("http://localhost:8000/sale", "application/json", body); err != nil {
		log.Panic(err)
	}

}

func sendMiembro(m miembro.Miembro) {

	txt, err := m.JSON()
	if err != nil {
		log.Panic(err)
	}

	body := strings.NewReader(string(txt))

	if _, err = http.Post("http://localhost:8000/member", "application/json", body); err != nil {
		log.Panic(err)
	}
}

func sendCoords(c coordinates.Coordinates) {

	txt, err := c.JSON()
	if err != nil {
		log.Panic(err)
	}

	body := strings.NewReader(string(txt))

	if _, err = http.Post("http://localhost:8000/coords", "application/json", body); err != nil {
		log.Panic(err)
	}
}

func getMember() miembro.Miembro {
	var (
		nombre   string
		apellido string
		rut      string
		email    string
		patente  string
		premium  bool
	)

	fmt.Print("Ingrese el nombre del miembro: ")
	fmt.Scanf("%s", &nombre)

	fmt.Print("Ingrese el apellido del miembro: ")
	fmt.Scanf("%s", &apellido)

	fmt.Print("Ingrese el rut del miembro: ")
	fmt.Scanln(&rut)

	for len(rut) > 12 {
		fmt.Println("El rut debe ser menor o igual a 12 caracteres.")
		fmt.Print("Ingrese nuevamente: ")
		fmt.Scanln(&rut)
	}

	fmt.Print("Ingrese el email del miembro: ")
	fmt.Scanf("%s", &email)

	fmt.Print("Ingrese la patente del miembro: ")
	fmt.Scan(&patente)

	for len(patente) > 6 {
		fmt.Println("La patente debe tener 6 o menos caracteres")
		fmt.Print("Ingrese nuevamente: ")
		fmt.Scan(&patente)
	}

	fmt.Print("Es usuario premium?: ")
	fmt.Scan(&premium)

	currMember = rut

	return miembro.Miembro{
		Nombre:   nombre,
		Apellido: apellido,
		Rut:      rut,
		Email:    email,
		Patente:  patente,
		Premium:  premium,
	}
}

func getVenta() (*venta.Venta, bool) {

	if currMember == "" {
		return nil, false
	}

	var (
		cliente  string
		cantidad int
		hora     string
		stock    int
		x, y     float64
	)

	fmt.Print("Ingrese el nombre del cliente: ")
	fmt.Scan(&cliente)

	fmt.Print("Ingrese la cantidad: ")
	fmt.Scan(&cantidad)

	fmt.Print("Ingrese la hora: ")
	fmt.Scan(&hora)

	fmt.Print("Ingrese el stock restante: ")
	fmt.Scan(&stock)

	fmt.Print("Ingrese las coordenadas de la forma 'x y': ")
	fmt.Scan(&x, &y)

	return &venta.Venta{
		Maestro:  currMember,
		Cliente:  cliente,
		Cantidad: cantidad,
		Hora:     hora,
		Stock:    stock,
		Coords:   [2]float64{x, y},
	}, true
}

func getCoords() coordinates.Coordinates {
	var x, y float64

	fmt.Print("Ingrese las coordenadas de la forma 'x y': ")
	fmt.Scan(&x, &y)

	return coordinates.Coordinates{
		Coords: [2]float64{x, y},
	}
}

func getSesion() {
	var rut string

	fmt.Print("Ingrese su rut: ")
	fmt.Scan(&rut)

	for len(rut) > 12 {
		fmt.Println("El largo del rut debe ser de 12 o menos caracteres.")
		fmt.Print("Ingrese nuevamente: ")
		fmt.Scan(&rut)
	}

	currMember = rut
}

func main() {
	var opcion int

	go sendLocation()

	for opcion != 5 {

		fmt.Printf("----------MENU----------\n\n")
		fmt.Println("1. Iniciar sesion")
		fmt.Println("2. Ingresar miembro ")
		fmt.Println("3. Enviar venta ")
		fmt.Println("4. Enviar coordenadas carrito profugo")
		fmt.Println("5. Salir")
		fmt.Println("\n--------INGRESE OPCION: ")

		fmt.Scanf("%d", &opcion)

		switch opcion {
		case 1:
			getSesion()
		case 2:
			sendMiembro(getMember())
		case 3:
			v, valid := getVenta()

			if valid {
				sendVentas(*v)
			} else {
				fmt.Println("No puede ingresar una venta sin iniciar sesion")
			}
		case 4:
			sendCoords(getCoords())
		}

	}

}
