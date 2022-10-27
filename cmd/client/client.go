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

var member []string
var r *rand.Rand
var n int

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
		xy := generateCoords(0, 100)
		c := coordinates.Coordinates{
			Coords:  xy,
			Miembro: member[n],
		}
		sendCoords(c)
		<-ticker.C
	}
}

func sendVentas(v venta.Venta) {

	v.Maestro = member[n]

	txt, err := v.JSON()
	if err != nil {
		log.Panic(err)
	}

	body := strings.NewReader(string(txt))

	_, err = http.Post("http://localhost:8000/sale", "application/json", body)
	if err != nil {
		log.Panic(err)
	}
	return

}

func sendMiembro(m miembro.Miembro) {

	m.Rut = member[n]

	txt, err := m.JSON()
	if err != nil {
		log.Panic(err)
	}

	body := strings.NewReader(string(txt))

	http.Post("http://localhost:8000/member", "application/json", body)
	return
}

func sendCoords(c coordinates.Coordinates) {

	txt, err := c.JSON()
	if err != nil {
		log.Panic(err)
	}

	body := strings.NewReader(string(txt))

	http.Post("http://localhost:8000/coords", "application/json", body)
	return
}

func main() {

	member = []string{"1234567-8", "3456789-0", "9876543-2"}
	r = rand.New(rand.NewSource(time.Now().Unix()))
	n = r.Intn(len(member))

	//EJEMPLOS
	m := miembro.Miembro{
		Nombre:   "juanito",
		Apellido: "Perez",
		Rut:      member[n],
		Email:    "juantio42069@hotmail.com",
		Patente:  "jdh3et23",
		Premium:  false,
	}

	v := venta.Venta{
		Maestro:  member[n],
		Cliente:  "wo",
		Cantidad: 10,
		Hora:     "13:52",
		Stock:    32,
		Coords:   [2]float64{28, 79},
	}

	c := coordinates.Coordinates{
		Coords: [2]float64{50, 136},
	}
	//EJEMPLOS

	fmt.Printf("El miembro %s ha iniciado sesion.\n", member[n])
	var opcion int

	go sendLocation()

	for opcion != 4 {

		fmt.Println("----------MENU----------\n")
		fmt.Println("1. Ingresar miembro ")
		fmt.Println("2. Enviar venta ")
		fmt.Println("3. Enviar coordenadas carrito profugo")
		fmt.Println("4. Salir")
		fmt.Println("\n--------INGRESE OPCION: ")

		fmt.Scanf("%d", &opcion)

		switch opcion {
		case 1:
			sendMiembro(m)
			break
		case 2:
			sendVentas(v)
			break
		case 3:
			sendCoords(c)
			break
		}

	}

}
