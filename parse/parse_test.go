package parse_test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/hoglandets-it/go-bankgiro/parse"
	"github.com/hoglandets-it/go-bankgiro/tools"
)

func TestAgainstExpected(t *testing.T) {
	fsd := os.DirFS("/mnt/f/BankFiles/AutogiroInbound/")

	files, err := fs.ReadDir(fsd, ".")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		t.Log("Testing " + file.Name())

		content, err := os.ReadFile("/mnt/f/BankFiles/AutogiroInbound/" + file.Name())
		if err != nil {
			t.Fatal(err)
		}

		isoContent, err := tools.BytesToIsoString(content)
		if err != nil {
			t.Error(err)
		}

		agFile := parse.AutogiroFile{}
		err = agFile.ParseFile(isoContent)
		if err != nil {
			t.Error(err)
		}

		marshalled, err := json.MarshalIndent(agFile, "", "  ")
		if err != nil {
			t.Error(err)
		}

		fmt.Println(string(marshalled))

		for i, section := range agFile.Sections {
			fmt.Println("sectionType: " + section.SectionType.Name)
			fmt.Println("customerNo: " + section.GetCustomerNumber())
			fmt.Println("accountNo: " + section.GetAccountNumber())

			utf8Bytes := section.GetUtf8Bytes()

			os.WriteFile(fmt.Sprintf("/mnt/f/BankFiles/AutogiroInbound/%s-section%d.txt", file, i), utf8Bytes, 0644)
		}

		fmt.Println("Done")
	}

}
