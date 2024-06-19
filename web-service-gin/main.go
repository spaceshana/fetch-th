package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"regexp"

	"github.com/google/uuid"

	"math"

	"strings"

	"strconv"

	"time"


)

// receipt represents data about a transaction at a retailer.
type receipt struct {
		Retailer string `json:"retailer"`
		PurchaseDate string `json:"purchaseDate"`
		PurchaseTime string `json:"purchaseTime"`
		Items []struct {
			ShortDescription string `json:"shortDescription"`
			Price string `json:"price"`
		} `json:"items"`
		Total string `json:"total"` 
}

type transaction struct {
	ID string `json:"id"`
	Points int `json:"points"`
}

type id struct {
	ID string `json:"id"`
}

type points struct {
	Points int `json:"points"`
}

func main() {
	router := gin.Default()
	router.POST("/receipts/process", postReceipt)
	router.GET("/receipts/:id/points", getReceiptPointsByID)

	router.Run("localhost:8080")

}


func postReceipt(c *gin.Context){
	var newReceipt receipt
	var receiptTransaction transaction
	var receiptId id
	if err := c.BindJSON(&newReceipt); err!=nil {
		c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "The receipt is invalid"})
		return
	}


	nums := process(newReceipt)

	receiptId.ID = uuid.New().String()
	
	receiptTransaction.ID = receiptId.ID
	receiptTransaction.Points = nums
	
	t = append(t, receiptTransaction)

	
	c.IndentedJSON(http.StatusOK, receiptId)
}

func process(r receipt) int{
	var p1,p2,p3,p4,p5,p6,p7 int

	total := r.Total
	items := r.Items
	purchaseDate := r.PurchaseDate
	purchaseTime := r.PurchaseTime
	purchaseDateTime := purchaseDate + " " + purchaseTime

	// 1 point for every alphanumeric character in the retailer name
	for _,ch := range r.Retailer{
		if(regexp.MustCompile("^[a-zA-Z0-9]*$").MatchString(string(ch))){
			p1 += 1
		}
	}

	// 50 points if the total is a round dollar amount with no cents
	if(total[len(total)-2:]=="00"){
		p2 = 50
	}

	// 25 points if the total is a multiple of 0.25
	if(math.Mod(convertStrtoNum(total),0.25)==0){
		p3 = 25
	}

	// 5 points for every two items on the receipt
	evenItems := len(items)/2
	p4 = 5*evenItems

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _,i := range items{
		itemDesc := strings.Trim(i.ShortDescription," ")
		
		if(len(itemDesc)%3==0){
			p5 += int(math.Ceil(convertStrtoNum(i.Price)*0.2))
		}

	}

	// 6 points if the day in the purchase date is odd.
	if(int(convertStrtoNum(purchaseDate[len(purchaseDate)-2:]))%2!=0){
		p6 = 6
	}
	
	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseYear := convertStrtoNum(purchaseDate[0:4])
	purchaseMonth := convertStrtoNum(purchaseDate[5:7])
	purchaseDay := convertStrtoNum(purchaseDate[8:10])
 
	after2 := time.Date(int(purchaseYear),time.Month(purchaseMonth),int(purchaseDay),14,00,00,00,time.UTC)
	before4 := time.Date(int(purchaseYear),time.Month(purchaseMonth),int(purchaseDay),16,00,00,00,time.UTC)

	if pt, err := time.Parse("2006-01-02 15:04", purchaseDateTime); err == nil{
		if(pt.After(after2) && pt.Before(before4)){
			p7 = 10
		}
	}

	sum := p1 + p2 + p3 + p4 + p5 + p6 + p7
	return sum
}

func convertStrtoNum(input string) float64{
	if s,err :=strconv.ParseFloat(input,32); err == nil{
		return s
	}

	return 0
}

func getReceiptPointsByID(c *gin.Context){
	var p points
	id := c.Param("id")

	for _, a := range t {
		if a.ID == id {
			p.Points = a.Points
			c.IndentedJSON(http.StatusOK, p)
			return
		}
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No receipt found for that id"})
}

var t =[]transaction{
}

