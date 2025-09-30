package printProduct

import (
	"WebCrauler/internal/Product"
	"fmt"

	"github.com/cheynewallace/tabby"
)

func PrintProducts(products []*product.Product, serchProduct string) {
	t := tabby.New()

	fmt.Printf("Serching product \"%s\" finaled: \n", serchProduct)

	t.AddHeader("Num", "ID", "IMG", "Price", "Title")
	for num, pr := range products {
		t.AddLine(num+1, pr.ID, pr.Img, pr.Price, pr.Title)
	}
	t.Print()
}