package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Item struct {
	Name     string
	Quantity int
	Price    float64
}

type Inventory struct {
	items map[string]Item
	file  string
}

func NewInventory(filename string) *Inventory {
	return &Inventory{
		items: make(map[string]Item),
		file:  filename,
	}
}

func (inv *Inventory) LoadFromFile() error {
	data, err := os.ReadFile(inv.file)
	if err != nil {
		if os.IsNotExist(err) {
			inv.items = make(map[string]Item)
			return nil
		}
		return err
	}

	// If the file is empty initialize an empty map
	if len(data) == 0 {
		inv.items = make(map[string]Item)
		return nil
	}

	// Decodes the json data into Inventory struct
	return json.Unmarshal(data, &inv.items)
}

func (inv *Inventory) AddItem(name string, quantity int, price float64) error {
	_, exists := inv.items[name]
	if exists {
		return fmt.Errorf("Item %s already present", name)
	}
	inv.items[name] = Item{
		Name:     name,
		Quantity: quantity,
		Price:    price,
	}
	return nil
}

func (inv *Inventory) UpdateQuantity(name string, addedQuant int) error {
	if addedQuant <= 0 {
		return fmt.Errorf("quantity added should be greater than 0")
	}
	item, exists := inv.items[name]
	if !exists {
		return fmt.Errorf("Item %s not found", name)
	}
	existingQunat := item.Quantity
	updatedQuant := existingQunat + addedQuant
	item.Quantity = updatedQuant
	inv.items[name] = item
	return nil
}

func (inv *Inventory) SaveToFile() error {
	// Marshal converts struct to json format
	data, err := json.MarshalIndent(inv.items, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}
	fmt.Printf("Saving data: %s\n", string(data))

	/*
		First digit (6): Owner can read (4) and write (2)
		Second digit (4): Group can read (4)
		Third digit (4): Others can read (4)
	*/
	err = os.WriteFile(inv.file, data, 0666)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}

func (inv *Inventory) DisplayItems() error {
	for _, item := range inv.items {
		fmt.Printf("item: %v, quantity: %v, price: %v \n", item.Name, item.Quantity, item.Price)
	}
	return nil
}

func main() {
	inventory := NewInventory("inventory.json")
	err := inventory.LoadFromFile()
	if err != nil {
		fmt.Println("Error creating/reading the inventory file")
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("1. Add new item")
		fmt.Println("2. Update quantity")
		fmt.Println("3. Display inventory")
		fmt.Println("4. Save the file")
		fmt.Println("5. Exit")

		fmt.Println("Choose an option")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			fmt.Println("Enter item name")
			scanner.Scan()
			name := scanner.Text()

			fmt.Println("Enter quantity")
			scanner.Scan()
			quantity, _ := strconv.Atoi(scanner.Text())

			fmt.Println("Enter price")
			scanner.Scan()
			price, _ := strconv.ParseFloat(scanner.Text(), 64)

			err := inventory.AddItem(name, quantity, price)
			if err != nil {
				fmt.Printf("Error: %v \n", err)
			} else {
				fmt.Println("Item added successfully")
			}

		case "2":
			fmt.Println("Enter item name")
			scanner.Scan()
			name := scanner.Text()

			fmt.Println("Enter quantity")
			scanner.Scan()
			quantityAdded, _ := strconv.Atoi(scanner.Text())
			err := inventory.UpdateQuantity(name, quantityAdded)
			if err != nil {
				fmt.Printf("Error: %v \n", err)
			} else {
				fmt.Println("Quantity updated successfully")
			}

		case "3":
			inventory.DisplayItems()

		case "4":
			err := inventory.SaveToFile()
			if err != nil {
				fmt.Printf("Error: %v \n", err)
			}
			fmt.Println("File saved successfully")

		case "5":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}
