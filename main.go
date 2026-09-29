package main

import (
	"encoding/json"
	"errors"
	"image/color"
	"os"
	"passwordEncryption/encryption"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var APP fyne.App = app.New()
var WINDOW fyne.Window = APP.NewWindow("Cipher Vault")
var ISCONNECTED bool = false
var PASSKEY string = ""
var VAULTPATH string = ""
var DATA = map[string]map[string]passwordsData{}
var USECASE = ""
var USERNAME = ""
var CURRENTPASSWORD = ""
var OLDPASSWORD = ""
var PAGETYPE pageType

type pageType int

const (
	HOME pageType = iota
	NEWFLE
	CONNECT
	ADD
	REMOVE
	EDIT
)

type passwordsData struct {
	current string
	histroy []string
}

func createHeader() (fyne.CanvasObject, error) {
	homeImg, err := fyne.LoadResourceFromPath("img/home100.png")
	if err != nil {
		return nil, errors.New("wasn't able to load home image")
	}
	newFileImg, err := fyne.LoadResourceFromPath("img/newFile100.png")
	if err != nil {
		return nil, errors.New("wasn't able to load newFile image")
	}
	connectImg, err := fyne.LoadResourceFromPath("img/connect100.png")
	if err != nil {
		return nil, errors.New("wasn't able to load connect image")
	}
	addImg, err := fyne.LoadResourceFromPath("img/add100.png")
	if err != nil {
		return nil, errors.New("wasn't able to load add image")
	}
	removeImg, err := fyne.LoadResourceFromPath("img/remove100.png")
	if err != nil {
		return nil, errors.New("wasn't able to load remove image")
	}
	editImg, err := fyne.LoadResourceFromPath("img/edit100.png")
	if err != nil {
		return nil, errors.New("wasn't able to load edit image")
	}

	homeButton := widget.NewButtonWithIcon("Home", homeImg, func() {
		println("Home button clicked")
	})
	newFileButton := widget.NewButtonWithIcon("New File", newFileImg, func() {
		println("New File button clicked")
	})
	connectButton := widget.NewButtonWithIcon("Connect", connectImg, func() {
		println("Connect button clicked")
	})
	addButton := widget.NewButtonWithIcon("Add", addImg, func() {
		println("Add button clicked")
	})
	removeButton := widget.NewButtonWithIcon("Remove", removeImg, func() {
		println("Remove button clicked")
	})
	editButton := widget.NewButtonWithIcon("Edit", editImg, func() {
		println("edit button clicked")
	})
	emptySpace := widget.NewLabel("")
	emptySpace.Resize(fyne.NewSize(30, 0))
	paddingSize := canvas.NewRectangle(color.Transparent)
	paddingSize.SetMinSize(fyne.NewSize(10, 10))
	buttons := container.NewBorder(paddingSize, paddingSize, paddingSize, paddingSize, container.NewHBox(homeButton, emptySpace, newFileButton, emptySpace, connectButton, emptySpace, addButton, emptySpace, removeButton, emptySpace, editButton))

	return container.NewVBox(buttons, canvas.NewLine(color.RGBA{R: 135, G: 135, B: 135, A: 255})), nil
}

func createUseCasesBox() fyne.CanvasObject {
	searchEntry := widget.NewEntry()
	searchEntry.PlaceHolder = "search your use case"
	searchButton := widget.NewButton("search", func() {

	})
	searchBox := container.NewBorder(
		nil,
		nil,
		nil,
		searchButton,
		searchEntry,
	)

	userNamesBox := container.NewStack()
	useCasesButtons := make([]fyne.CanvasObject, 0, 20)
	for useCase, _ := range DATA {
		useCasesButtons = append(useCasesButtons, widget.NewButton(useCase, func() {
			userNamesBox.Objects = []fyne.CanvasObject{createUserNamesBox(useCase)}
			userNamesBox.Refresh()
			USECASE = useCase
		}))
	}

	useCasesBox := container.NewBorder(
		searchBox,
		nil,
		nil,
		nil,
		container.NewScroll(container.NewVBox(useCasesButtons...)),
	)

	result := container.NewHSplit(useCasesBox, userNamesBox)
	result.SetOffset(0.4)
	return result
}

func createUserNamesBox(useCase string) fyne.CanvasObject {

	passwordBox := container.NewStack()
	userNameButtons := make([]fyne.CanvasObject, 0, 5)
	for userName, passwordsData := range DATA[useCase] {
		userNameButtons = append(userNameButtons, widget.NewButton(userName, func() {
			USERNAME = userName
			passwordBox.Objects = []fyne.CanvasObject{createPasswordsBox(passwordsData)}
			passwordBox.Refresh()
		}))
	}
	userNamesBox := container.NewBorder(
		widget.NewLabelWithStyle("UserNames:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewScroll(
			container.NewVBox(userNameButtons...),
		),
	)
	userNameAndPasswordSplitter := container.NewHSplit(userNamesBox, passwordBox)
	userNameAndPasswordSplitter.SetOffset(0.4)
	return userNameAndPasswordSplitter
}

func createPasswordsBox(passwordsData passwordsData) fyne.CanvasObject {
	oldPasswords := []fyne.CanvasObject{}
	for _, password := range passwordsData.histroy {
		oldPasswords = append(oldPasswords, createPasswordBox(password))
	}

	return container.NewBorder(
		container.NewVBox(widget.NewLabelWithStyle("Password:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			createPasswordBox(passwordsData.current),
			widget.NewLabelWithStyle("History:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		nil,
		nil,
		nil,
		container.NewScroll(
			container.NewVBox(oldPasswords...),
		),
	)
}

func createPasswordBox(password string) fyne.CanvasObject {
	compImg, _ := fyne.LoadResourceFromPath("img/copy.png")
	return container.NewStack(
		canvas.NewRectangle(color.RGBA{R: 43, G: 55, B: 69, A: 255}),
		container.NewBorder(
			nil,
			nil,
			nil,
			container.NewPadded(widget.NewButtonWithIcon("", compImg, func() {
				WINDOW.Clipboard().SetContent(password)
			})),
			widget.NewLabelWithStyle(password, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		),
	)
}

func readFile() error {

	encrypter := &encryption.Encrypter{PassKey: PASSKEY}
	rawData, err := os.ReadFile(VAULTPATH)
	if err != nil {
		return errors.New("wasn't able to read file from " + VAULTPATH)
	}

	data, err := encrypter.Decrypt(rawData)
	if err != nil {
		return errors.New("wasn't able to read file from " + VAULTPATH)
	}

	if err = json.Unmarshal(data, &DATA); err != nil {
		return errors.New("couldn't unmarshal json file")
	}

	for useCase := range DATA {
		for _, passwordData := range DATA[useCase] {
			currentPassword, err := encrypter.Decrypt([]byte(passwordData.current))
			if err != nil {
				return errors.New("wrong PassKey")
			}
			passwordData.current = string(currentPassword)
			for i, password := range passwordData.histroy {
				oldPassword, err := encrypter.Decrypt([]byte(password))
				if err != nil {
					return errors.New("wrong PassKey")
				}
				passwordData.histroy[i] = string(oldPassword)
			}
		}
	}
	return nil
}

func createContainer() fyne.CanvasObject {
	var ctr fyne.CanvasObject
	if !ISCONNECTED {
		ctr = container.NewCenter(
			container.NewStack(
				canvas.NewRectangle(color.RGBA{R: 156, G: 0, B: 54, A: 255}),
				widget.NewLabelWithStyle("NOTE CONNECTED TO ANY VAULT", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			),
		)
		return ctr
	}
	switch PAGETYPE {
	case HOME:
		ctr = createUseCasesBox()
	}
	return ctr
}

func createFooter() fyne.CanvasObject {
	footerText := container.NewStack()
	tempMessage := ""
	if ISCONNECTED {
		tempMessage = "Connected to " + VAULTPATH
	} else {
		tempMessage = "Not connected to any vault"
	}
	rectangle := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 20})
	text := canvas.NewText(tempMessage, color.White)
	text.TextSize = 13
	footerText.Objects = []fyne.CanvasObject{
		rectangle,
		container.NewPadded(text),
	}
	return footerText
}

func createNewEncrypted() error {
	file, err := os.Create(VAULTPATH)
	if err != nil {
		return errors.New("couldn't create the file")
	}

	defer file.Close()
	return nil
}

// func createNewFileBox() fyne.CanvasObject {

// }

func main() {
	bgColor := canvas.NewRectangle(color.RGBA{R: 25, G: 25, B: 26, A: 255})
	WINDOW.Resize(fyne.NewSize(750, 450))

	// w.SetFixedSize(true)
	h, err := createHeader()
	if err != nil {
		panic(err)
	}

	DATA = map[string]map[string]passwordsData{
		"udemy": {
			"Sowhail":  {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"Sowbham":  {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"gham":     {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"irjfg":    {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"gam":      {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"Surhj3":   {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"fhjvjkm":  {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"vgvcm":    {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"irhjham":  {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"Sowbfj":   {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"fgjm":     {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"ifbham":   {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
			"fbvwbham": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"udgemy": {
			"sobha": {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"wothg": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"sdf": {
			"lbvch":   {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"iuwehgn": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"hedh": {
			"wowow": {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"lxcvh": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"hegdh": {
			"wowow": {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"lxcvh": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"hefdh": {
			"wowow": {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"lxcvh": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"hedagh": {
			"wowow": {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"lxcvh": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"hedhz": {
			"wowow": {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"lxcvh": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"hedbh": {
			"wowow": {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"lxcvh": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
		"hedwh": {
			"wowow": {current: "sl3gfhfd", histroy: []string{"dalright", "hwe52", "ghfo0"}},
			"lxcvh": {current: "qoglb", histroy: []string{"jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg", "jfdgh", "alre452", "g82hg"}},
		},
	}

	window := container.NewStack(bgColor, container.NewBorder(
		h,
		createFooter(),
		nil,
		nil,
		createContainer(),
	))

	WINDOW.SetContent(window)

	// header Buttons

	WINDOW.Show()
	APP.Run()
}
