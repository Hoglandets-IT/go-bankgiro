package parse

import (
	"fmt"
	"slices"
	"strings"
)

type SectionInterface struct {
	GetAccountNumber  func() string
	GetCustomerNumber func() string
	GetUtf8Bytes      func() []byte
	GetUtf8String     func() []byte
}

type SectionType struct {
	Name            string
	Code            string
	Tk01Start       int
	Tk01End         int
	Match           string
	AllowedSections []string
	CustomerNumber  []int
	AccountNumber   []int
}

const (
	HMAC_HEADER       = "00"
	SECTION_START     = "01"
	SECTION_END       = "09"
	HMAC_SECTION_SEAL = "08"
	HMAC_FILE_SEAL    = "99"
)

var SectionTypes []SectionType = []SectionType{
	{ // OK
		Name:            "Utdrag Bevakningsregister (Gammalt/Nytt Format)",
		Code:            "bevakningsreg",
		Tk01Start:       22,
		Tk01End:         35,
		Match:           "BEVAKNINGSREG",
		AllowedSections: []string{"82", "32"},
		CustomerNumber:  []int{62, 68},
		AccountNumber:   []int{68, 78},
	},
	{ // OK
		Name:            "Medgivandeavisering (Gammalt Format)",
		Code:            "medgivandeavi-old",
		Tk01Start:       24,
		Tk01End:         33,
		Match:           "AG-MEDAVI",
		AllowedSections: []string{"73"},
		CustomerNumber:  []int{0, 0},
		AccountNumber:   []int{14, 24},
	},
	{ // OK
		Name:            "Avvisade Betalningar (Gammalt Format)",
		Code:            "avvisade-old",
		Tk01Start:       22,
		Tk01End:         41,
		Match:           "FELLISTA REG.KONTRL",
		AllowedSections: []string{"82", "32"},
		CustomerNumber:  []int{62, 68},
		AccountNumber:   []int{68, 78},
	},
	{ // OK
		Name:            "Makulerings-/Ändringslista (Gammalt Format)",
		Code:            "andringslista-old",
		Tk01Start:       22,
		Tk01End:         40,
		Match:           "MAK/ÄNDRINGSLISTA",
		AllowedSections: []string{"03", "21", "22", "23", "24", "25", "26", "27", "28", "29"},
		CustomerNumber:  []int{62, 68},
		AccountNumber:   []int{68, 78},
	},
	{ // OK
		Name:            "Betalningsspecifikation (Nytt Format)",
		Code:            "betalningsspec-new",
		Tk01Start:       44,
		Tk01End:         64,
		Match:           "BET. SPEC & STOPP TK",
		AllowedSections: []string{"15", "82", "16", "32", "17", "77"},
		CustomerNumber:  []int{64, 70},
		AccountNumber:   []int{70, 80},
	},
	{ // OK
		Name:            "Medgivandeavisering (Nytt Format)",
		Code:            "medgivandeavi-new",
		Tk01Start:       44,
		Tk01End:         53,
		Match:           "AG-MEDAVI",
		AllowedSections: []string{"73"},
		CustomerNumber:  []int{64, 70},
		AccountNumber:   []int{70, 80},
	},
	{ // OK
		Name:            "Avvisade Betalningar (Nytt Format)",
		Code:            "avvisade-new",
		Tk01Start:       44,
		Tk01End:         62,
		Match:           "AVVISADE BET UPPDR",
		AllowedSections: []string{"82", "32"},
		CustomerNumber:  []int{64, 70},
		AccountNumber:   []int{70, 80},
	},
	{ // OK
		Name:            "Makulerings-/Ändringslista (Nytt Format)",
		Code:            "andringslista-new",
		Tk01Start:       44,
		Tk01End:         63,
		Match:           "MAKULERING/ÄNDRING",
		AllowedSections: []string{"03", "11", "21", "22", "23", "24", "25", "26", "27", "28", "29"},
		CustomerNumber:  []int{64, 70},
		AccountNumber:   []int{70, 80},
	},
	{
		Name:            "Betalningsspecifikation (Gammalt Format)",
		Code:            "betalningsspec-old",
		Tk01Start:       10,
		Tk01End:         18,
		Match:           "AUTOGIRO",
		AllowedSections: []string{"82", "32"},
		CustomerNumber:  []int{62, 68},
		AccountNumber:   []int{68, 78},
	},
	{
		Name:            "INVALID FILE TYPE",
		Code:            "invalid",
		Tk01Start:       0,
		Tk01End:         2,
		Match:           "01",
		AllowedSections: []string{},
		CustomerNumber:  []int{0, 0},
		AccountNumber:   []int{0, 0},
	},
}

type AutogiroSection struct {
	StartFound      bool
	SectionType     SectionType
	EndFound        bool
	SectionSeal     string
	Rows            []string
	SealCalcContent []string
	Errors          []string
}

type AutogiroFile struct {
	HMACStartFound  bool
	HMACEndFound    bool
	HmacData        string
	Content         []string
	SealCalcContent []string
	Sections        []AutogiroSection
}

func (sec *AutogiroSection) SetStart(line string) error {
	for _, sectionType := range SectionTypes {
		substr := line[sectionType.Tk01Start:sectionType.Tk01End]
		if substr == sectionType.Match {
			sec.StartFound = true
			sec.Rows = append(sec.Rows, line)
			sec.SectionType = sectionType
			return nil
		}
	}

	return fmt.Errorf("no matching section type found")
}

func (sec *AutogiroSection) SetEnd(line string, lookaheadRow string) error {
	sec.EndFound = true
	sec.Rows = append(sec.Rows, line)

	if len(lookaheadRow) > 1 && lookaheadRow[0:2] == HMAC_SECTION_SEAL {
		sec.SectionSeal = lookaheadRow
	}

	return nil
}

// Input: Line
// Output: End of Section
func (sec *AutogiroSection) AddLine(line string) error {
	if len(line) != 80 {
		sec.Errors = append(sec.Errors, fmt.Sprintf("Invalid line length: %d - *%s*", len(line), line))
	}
	sec.Rows = append(sec.Rows, line)
	if !slices.Contains(sec.SectionType.AllowedSections, line[0:2]) {
		sec.Errors = append(sec.Errors, fmt.Sprintf("Invalid section found: %s", line[0:2]))
	}

	return nil
}

func (sec *AutogiroSection) GetAccountNumber() string {
	return sec.Rows[0][sec.SectionType.AccountNumber[0]:sec.SectionType.AccountNumber[1]]
}

func (sec *AutogiroSection) GetCustomerNumber() string {

	return sec.Rows[0][sec.SectionType.CustomerNumber[0]:sec.SectionType.CustomerNumber[1]]
}

func (sec *AutogiroSection) GetUtf8Bytes() []byte {
	return []byte(strings.Join(sec.Rows, "\r\n"))
}

func (sec *AutogiroSection) GetUtf8String() string {
	return strings.Join(sec.Rows, "\r\n")
}

func (file *AutogiroFile) ParseFile(data string) error {
	file.HMACStartFound = false
	file.HMACEndFound = false

	file.Sections = make([]AutogiroSection, 0)

	// Split the file into rows
	file.Content = strings.Split(data, "\r\n")
	if len(file.Content) == 1 {
		file.Content = strings.Split(data, "\n")
		if len(file.Content) == 1 {
			file.Content = strings.Split(data, "\r")
			if len(file.Content) == 1 {
				return fmt.Errorf("could not split the file into rows: only one row present")
			}
		}
	}

	// Loop through the rows and parse the file
	var currentSection AutogiroSection = AutogiroSection{}

	for i, row := range file.Content {
		// Identify new section
		if currentSection.StartFound && currentSection.EndFound {
			file.Sections = append(file.Sections, currentSection)
			currentSection = AutogiroSection{}
		}

		// Identify empty rows
		if strings.Trim(row, " \t") == "" || row[0:2] == HMAC_SECTION_SEAL || row[0:2] == HMAC_FILE_SEAL {
			continue
		}

		// Identify HMAC Start row
		if strings.HasPrefix(row, HMAC_HEADER) {

			// If start was already found, throw an error
			if file.HMACStartFound {
				return fmt.Errorf("multiple hmac start lines found")
			}

			file.HMACStartFound = true
			continue
		}

		// Identify Section Starter
		if !currentSection.StartFound {
			if !strings.HasPrefix(row, SECTION_START) {
				return fmt.Errorf("no section start found where there should be one")
			}

			err := currentSection.SetStart(row)
			if err != nil {
				return err
			}

			continue
		}

		// Identify section end
		if currentSection.StartFound && !currentSection.EndFound {
			if strings.HasPrefix(row, SECTION_END) {
				lookaheadRow := ""
				if len(file.Content) >= i+2 {
					lookaheadRow = file.Content[i+1]
				}

				err := currentSection.SetEnd(row, lookaheadRow)
				if err != nil {
					return err
				}

				file.Sections = append(file.Sections, currentSection)
				currentSection = AutogiroSection{}

				continue
			}
		}

		err := currentSection.AddLine(row)
		if err != nil {
			return err
		}
	}

	if file.HMACStartFound && !file.HMACEndFound {
		if file.Content[len(file.Content)-1][0:2] != HMAC_FILE_SEAL {
			return fmt.Errorf("hmac never ended")
		}

		file.HMACEndFound = true
		file.HmacData = file.Content[len(file.Content)-1]
	}

	if currentSection.StartFound && !currentSection.EndFound {
		file.Sections = append(file.Sections, currentSection)
		return fmt.Errorf("section never ended")
	}

	return nil
}
