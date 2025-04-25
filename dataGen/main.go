package main

import (
	"encoding/csv"
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
)

func pickFromSlice[T any](r *rand.Rand, slice []T) T {
	idx := r.Intn(len(slice))
	return slice[idx]
}

type User struct {
	Firstname string
	Lastname  string
	Age       uint
	Gender    string
}

func (self User) toCSV() []string {
	b := []string{}

	b = append(b, self.Firstname)
	b = append(b, self.Lastname)
	b = append(b, strconv.FormatUint(uint64(self.Age), 10))
	b = append(b, self.Gender)

	return b
}

func createRandomUser(r *rand.Rand) User {
	firstnames := []string{
		"Maria",
		"Jose",
		"Flavio",
		"Pablo",
		"Nicolle",
		"Julia",
	}

	lastnames := []string{
		"Galan",
		"Paz",
		"Guillermo",
		"Toc",
	}

	gender := []string{
		"Masculino",
		"Femenino",
	}

	return User{
		Firstname: pickFromSlice(r, firstnames),
		Lastname:  pickFromSlice(r, lastnames),
		Age:       uint(r.Intn(40) + 18),
		Gender:    pickFromSlice(r, gender),
	}
}

type MenuItem struct {
	Name  string
	Price float64
}

func createRandomMenuItem(r *rand.Rand) MenuItem {
	menuItemNames := []string{
		"Nigiri",
		"Miso Soup",
		"Pizza Pepperoni",
		"Pizza Margarita",
		"Butter Chicken",
	}

	name := pickFromSlice(r, menuItemNames)
	price := r.NormFloat64()

	return MenuItem{
		Name:  name,
		Price: price,
	}
}

type Menu struct {
	Items []MenuItem
}

func createRandomMenu(r *rand.Rand, itemCount uint) Menu {
	items := []MenuItem{}

	for range itemCount {
		items = append(items, createRandomMenuItem(r))
	}

	return Menu{Items: items}
}

type Location struct {
	LocType     string
	Coordinates []float64
}

func createRandomLocation(r *rand.Rand) Location {
	cords := []float64{
		r.NormFloat64(),
		r.NormFloat64(),
	}

	return Location{
		LocType:     "Point",
		Coordinates: cords,
	}
}

type Restaurant struct {
	name string
	// GeoJSON
	location Location
	// Date
	dob      string
	category string
	pricing  uint
	// ObjectId<file>
	photo string

	menu Menu
}

func (self Restaurant) toCSV() []string {
	b := []string{}

	b = append(b, self.name)

	location, err := json.Marshal(self.location)
	if err != nil {
		panic(err)
	}
	b = append(b, string(location))

	b = append(b, self.dob)
	b = append(b, self.category)
	b = append(b, strconv.FormatUint(uint64(self.pricing), 10))
	b = append(b, self.photo)

	menu, err := json.Marshal(self.menu)
	if err != nil {
		panic(err)
	}
	b = append(b, string(menu))

	return b
}

func createRandomRestaurant(r *rand.Rand) Restaurant {
	restaurantNames := []string{
		"Giorno",
		"Taco Fiesta",
		"Sushi Zen",
		"Pizza Mania",
	}
	categoryNames := []string{
		"Mexican",
		"Japanese",
		"American",
	}
	dobs := []string{
		"2015-06-12T00:00:00Z",
		"2018-03-22T00:00:00Z",
		"2020-01-15T00:00:00Z",
	}

	name := pickFromSlice(r, restaurantNames)
	location := createRandomLocation(r)
	dob := pickFromSlice(r, dobs)
	category := pickFromSlice(r, categoryNames)
	pricing := uint(r.Intn(4) + 1)
	photo := "1"
	menu := createRandomMenu(r, uint(r.Intn(10)))

	return Restaurant{
		name, location, dob, category, pricing, photo, menu,
	}
}

type csvable interface {
	toCSV() []string
}

func GenAndWriteToCSV[T csvable](r *rand.Rand, filename string, quantity uint, generator func(*rand.Rand) T, header []string) {
	data := [][]string{}
	for range quantity {
		rowData := generator(r)
		data = append(data, rowData.toCSV())
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(file)

	err = writer.Write(header)
	if err != nil {
		panic(err)
	}

	err = writer.WriteAll(data)
	if err != nil {
		panic(err)
	}
}

func main() {

	s := rand.NewSource(int64(696969))
	r := rand.New(s)

	os.RemoveAll("output")
	err := os.Mkdir("output", 0750)
	if err != nil {
		panic(err)
	}

	header := []string{
		"Name",
		"Location",
		"Dob",
		"Category",
		"Pricing",
		"Photo",
		"Menu",
	}
	GenAndWriteToCSV(r, "output/restaurants.csv", 50, createRandomRestaurant, header)

	header = []string{
		"Firstname",
		"Lastname",
		"Age",
		"Gender",
	}
	GenAndWriteToCSV(r, "output/users.csv", 300, createRandomUser, header)
}
