package main

import (
	"fmt"
	"wlog/editor"
	"wlog/store"
)

func main() {
	content, err := editor.Read()
	if err != nil {
		fmt.Println("error reading from editor:", err)
		return
	}

	fmt.Println("Content from editor:")
	fmt.Println(content)

	st := store.NewStore(store.Option{
		Directory: "./wlog_data",
	})

	err = st.SaveData([]byte(content))
	if err != nil {
		fmt.Println("error saving to store:", err)
		return
	}

	fmt.Println("Store saved.")

}
