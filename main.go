package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)


var mu sync.Mutex
var reportsList []string

func reports()  {
	mu.Lock()
	defer mu.Unlock()
	file, err := os.OpenFile("reports.txt",os.O_APPEND | os.O_CREATE |os.O_WRONLY,0644)
	if err != nil {
		fmt.Println("error in creating file")
		return
	}
	defer file.Close()

	for i,r := range reportsList{
		fmt.Fprintf(file,"%d. %s \n",i+1,r)
	}

	fmt.Println("reports added successfully")
}
func main()  {
	fmt.Println("server running on localhost : 8000")

	http.HandleFunc("/home",func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintln(w,"use post /weather for finding city weather information")
		} else {
			http.Error(w,"you should use get method",http.StatusMethodNotAllowed)
		}
		
	})

	http.HandleFunc("/weather",func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			city := r.FormValue("city")

			if city == "" {
				http.Error(w, "Missing city name", http.StatusBadRequest)
				return
			}

			report := fmt.Sprintf("the weather for %s is 27C",city)

			mu.Lock()
			reportsList = append(reportsList,report)
			mu.Unlock()

			fmt.Fprintln(w,report)
			reports()

		} else {
			http.Error(w,"only post allowed",http.StatusMethodNotAllowed)
		}
		
	})

	http.HandleFunc("/reports",func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

			mu.Lock()
			data,err := os.ReadFile("reports.txt")
			mu.Unlock()

			if err != nil {
				fmt.Fprintln(w,"error in reading file ")
				return
			} 
			fmt.Fprintln(w,string(data))
			
		} else {
			http.Error(w,"only get method allowed",http.StatusMethodNotAllowed)
		}
	})
	
	http.ListenAndServe(":8000",nil)
}