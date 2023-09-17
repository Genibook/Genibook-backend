package pages

import (
	constants "webscrapper/constants/v1"
	"webscrapper/utils"

	"github.com/gocolly/colly"
)

// func saveToFile(filename string, data []byte) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = file.Write(data)
// 	return err
// }

func GetPDFDataBytes(c *colly.Collector, school string, id string) ([]byte, error) {
	body := make([]byte, 0)

	c.OnResponse(func(r *colly.Response) {
		body = r.Body
		// pdfFilename := "downloaded.pdf"
		// err := saveToFile(pdfFilename, body)
		// if err != nil {
		// 	fmt.Println("Error saving PDF:", err)
		// 	return
		// }

		// fmt.Printf("PDF saved as %s\n", pdfFilename)
	})

	data := constants.ConstantLinks[school]["transcript"]
	data["studentid"] = id
	url, err := utils.FormatDynamicUrl(data, school)
	if err != nil {
		return body, err
	}

	err = c.Visit(url)
	if err != nil {
		return body, err
	}

	c.OnResponse(func(r *colly.Response) {})

	return body, nil
}
