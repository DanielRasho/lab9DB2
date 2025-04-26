package main

import (
	"encoding/csv"
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

const MAX_RESTAURANTS int = 50
const MAX_CLIENTS int = 400
const MAX_ORDERS int = 1000
const MAX_REVIEWS int = 1000

func pickFromSlice[T any](r *rand.Rand, slice []T) T {
	idx := r.Intn(len(slice))
	return slice[idx]
}

type Review struct {
	restaurant string
	client     string
	rating     uint
	relevance  uint
}

func (self Review) toCSV() []string {
	b := []string{}

	b = append(b, self.restaurant)
	b = append(b, self.client)
	b = append(b, strconv.FormatUint(uint64(self.rating), 10))
	b = append(b, strconv.FormatUint(uint64(self.relevance), 10))

	return b
}

func createRandomReview(r *rand.Rand) Review {

	restaurant := r.Intn(MAX_RESTAURANTS) + 1
	client := r.Intn(MAX_CLIENTS) + 1
	rating := r.Intn(11)
	relevance := r.Intn(5) + 1

	return Review{
		restaurant: strconv.FormatUint(uint64(restaurant), 10),
		client:     strconv.FormatUint(uint64(client), 10),
		rating:     uint(rating),
		relevance:  uint(relevance),
	}
}

type OrderItem struct {
	Name  string
	Price float64
}

func createRandomOrderItem(r *rand.Rand) OrderItem {
	menuItemNames := []string{
		"Nigiri",
		"Miso Soup",
		"Pizza Pepperoni",
		"Pizza Margarita",
		"Butter Chicken",
	}

	name := pickFromSlice(r, menuItemNames)
	price := r.Float64() * 150

	return OrderItem{
		Name:  name,
		Price: price,
	}

}

type Order struct {
	client     string
	restaurant string
	// : [“delivered”, “in-process”, “accepted”],
	state    string
	date     string
	pricing  float64
	quantity uint
	item     OrderItem
}

func (self Order) toCSV() []string {
	b := []string{}

	b = append(b, self.client)
	b = append(b, self.restaurant)
	b = append(b, self.state)
	b = append(b, self.date)
	b = append(b, strconv.FormatFloat(self.pricing, 'f', 2, 64))
	b = append(b, strconv.FormatUint(uint64(self.quantity), 10))

	item, err := json.Marshal(self.item)
	if err != nil {
		panic(err)
	}
	b = append(b, string(item))

	return b
}

func createRandomOrder(r *rand.Rand) Order {
	states := []string{
		"delivered", "in-process", "accepted",
	}

	dobs := []string{
		"2015-06-12T00:00:00Z",
		"2018-03-22T00:00:00Z",
		"2020-01-15T00:00:00Z",
	}

	client := r.Intn(MAX_CLIENTS) + 1
	restaurant := r.Intn(MAX_RESTAURANTS) + 1
	state := pickFromSlice(r, states)
	date := pickFromSlice(r, dobs)
	pricing := r.Float64() * 150
	quantity := r.Intn(10) + 1
	item := createRandomOrderItem(r)

	return Order{
		client:     strconv.FormatInt(int64(client), 10),
		restaurant: strconv.FormatInt(int64(restaurant), 10),
		state:      state,
		date:       date,
		pricing:    pricing,
		quantity:   uint(quantity),
		item:       item,
	}
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
	price := r.Float64() * 150

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
	csvCols := []string{"_id"}
	csvCols = append(csvCols, header...)
	data := [][]string{csvCols}

	for idx := range quantity {
		id := strconv.FormatUint(uint64(idx+1), 10)
		csvData := []string{id}

		rowData := generator(r)
		csvData = append(csvData, rowData.toCSV()...)

		data = append(data, csvData)
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(file)

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

	var group sync.WaitGroup

	group.Add(1)
	go func() {
		defer group.Done()
		header := []string{
			"Name",
			"Location",
			"Dob",
			"Category",
			"Pricing",
			"Photo",
			"Menu",
		}
		GenAndWriteToCSV(r, "output/restaurants.csv", uint(MAX_RESTAURANTS), createRandomRestaurant, header)
	}()

	group.Add(1)
	go func() {
		defer group.Done()
		header := []string{
			"Firstname",
			"Lastname",
			"Age",
			"Gender",
		}
		GenAndWriteToCSV(r, "output/users.csv", uint(MAX_CLIENTS), createRandomUser, header)
	}()

	group.Add(1)
	go func() {
		defer group.Done()
		header := []string{
			"Client",
			"Restaurant",
			"State",
			"Date",
			"Pricing",
			"Quantity",
			"Item",
		}
		GenAndWriteToCSV(r, "output/orders.csv", uint(MAX_ORDERS), createRandomOrder, header)
	}()

	group.Add(1)
	go func() {
		defer group.Done()
		header := []string{
			"Restaurant",
			"Client",
			"Rating",
			"Relevance",
		}
		GenAndWriteToCSV(r, "output/reviews.csv", uint(MAX_REVIEWS), createRandomReview, header)
	}()

	group.Wait()
}
